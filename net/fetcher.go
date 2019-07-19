package net

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/ledger"
	"github.com/koinotice/vite/log15"
)

var errNoSuitablePeer = errors.New("no suitable peer")

type gid struct {
	index uint32 // atomic
}

func (g *gid) MsgID() MsgId {
	return atomic.AddUint32(&g.index, 1)
}

type MsgIder interface {
	MsgID() MsgId
}

// fetch filter
const maxMark = 2       // times
const timeThreshold = 3 // second
const expiration = 20   // 20s

type record struct {
	id          MsgId
	addAt       int64
	t           int64 // change state time
	mark        int
	st          reqState
	targets     map[peerId]bool // whether failed
	restTargets int
	callback    func(msg Msg, err error)
}

func (r *record) inc() {
	r.mark++
}

func (r *record) clean() {
	r.restTargets = 0
	r.st = reqPending
	r.mark = 0
	r.addAt = time.Now().Unix()
}

func (r *record) reset() {
	r.mark = 0
	r.st = reqPending
	r.addAt = time.Now().Unix()
	r.targets = make(map[peerId]bool)
	r.restTargets = 0
	r.callback = nil
}

func (r *record) done(msg Msg, err error) {
	if r.st != reqPending {
		return
	}

	if err != nil {
		r.st = reqError
	} else {
		r.st = reqDone
	}
	r.t = time.Now().Unix()

	if r.callback != nil {
		go r.callback(msg, err)
	}
}

func (r *record) fail(pid peerId) {
	if failed, ok := r.targets[pid]; ok {
		if false == failed {
			r.targets[pid] = false
			r.restTargets--
		}
	}

	if r.restTargets == 0 && r.st == reqPending {
		r.st = reqError
		r.t = time.Now().Unix()
		if r.callback != nil {
			r.callback(Msg{}, errors.New("timeout"))
		}
	}
}

type filter struct {
	idGen    MsgIder
	idToHash map[MsgId]*record
	records  map[types.Hash]*record
	mu       sync.Mutex
	pool     sync.Pool
}

func newFilter() *filter {
	return &filter{
		idGen:    new(gid),
		idToHash: make(map[MsgId]*record, 1000),
		records:  make(map[types.Hash]*record, 1000),
		pool: sync.Pool{
			New: func() interface{} {
				return &record{
					targets: make(map[peerId]bool),
				}
			},
		},
	}
}

func (f *filter) clean(t int64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for hash, r := range f.records {
		if (t - r.addAt) > expiration {
			delete(f.records, hash)
			delete(f.idToHash, r.id)

			r.done(Msg{}, errors.New("timeout"))

			f.pool.Put(r)
		}
	}
}

// will suppress fetch
func (f *filter) hold(hash types.Hash) (r *record, hold bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now().Unix()

	var ok bool
	if r, ok = f.records[hash]; ok {
		if r.st == reqError {
			r.clean()
			return r, false
		} else if r.st == reqDone {
			if r.mark >= maxMark && (now-r.t) >= timeThreshold {
				r.clean()
				return r, false
			}
		} else {
			// pending
			if r.mark >= maxMark*2 && (now-r.addAt) >= timeThreshold*2 {
				r.clean()
				return r, false
			}
		}

		r.inc()
		return r, true
	}

	r = f.pool.Get().(*record)
	r.reset()
	r.id = f.idGen.MsgID()

	f.records[hash] = r
	f.idToHash[r.id] = r

	return r, false
}

func (f *filter) add(hash types.Hash) (r *record) {
	r = f.pool.Get().(*record)
	r.reset()
	r.id = f.idGen.MsgID()

	f.mu.Lock()
	f.idToHash[r.id] = r
	f.mu.Unlock()

	return r
}

func (f *filter) done(id MsgId, msg Msg) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r, ok := f.idToHash[id]; ok {
		r.done(msg, nil)
	}
}

func (f *filter) fail(id MsgId, pid peerId) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r, ok := f.idToHash[id]; ok {
		r.fail(pid)
	}
}

// height is 0 when fetch account block
func (f *filter) pickTargets(r *record, height uint64, peers *peerSet) peers {
	if height > 10 {
		height -= 10
	}

	ps := peers.pickReliable(height)

	if len(ps) == 0 {
		return nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	var hole bool
	for i, p := range ps {
		if _, ok := r.targets[p.Id]; ok {
			hole = true
			ps[i] = nil
			continue
		}
	}

	if hole {
		var j int
		for _, p := range ps {
			if p == nil {
				continue
			}
			ps[j] = p
			j++
		}

		ps = ps[:j]
	}

	if len(ps) > 3 {
		ps = ps[:3]
	}

	for _, p := range ps {
		r.targets[p.Id] = false
		r.restTargets++
	}

	return ps
}

type fetcher struct {
	filter *filter

	peers *peerSet

	st       SyncState
	receiver blockReceiver

	log log15.Logger

	blackBlocks map[types.Hash]struct{}
	sbp         bool

	term chan struct{}
}

func newFetcher(peers *peerSet, receiver blockReceiver, blackBlocks map[types.Hash]struct{}) *fetcher {
	if len(blackBlocks) == 0 {
		blackBlocks = make(map[types.Hash]struct{})
	}

	return &fetcher{
		filter:      newFilter(),
		peers:       peers,
		receiver:    receiver,
		blackBlocks: blackBlocks,
		log:         netLog.New("module", "fetcher"),
	}
}

func (f *fetcher) setSBP(bool2 bool) {
	f.sbp = bool2
}

func (f *fetcher) start() {
	f.term = make(chan struct{})
	go f.cleanLoop()
}

func (f *fetcher) stop() {
	if f.term == nil {
		return
	}

	select {
	case <-f.term:
	default:
		close(f.term)
	}
}

func (f *fetcher) cleanLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-f.term:
			return
		case now := <-ticker.C:
			f.filter.clean(now.Unix())
		}
	}
}

func (f *fetcher) subSyncState(st SyncState) {
	f.st = st
}

func (f *fetcher) name() string {
	return "fetcher"
}

func (f *fetcher) codes() []Code {
	return []Code{CodeSnapshotBlocks, CodeAccountBlocks, CodeException}
}

func (f *fetcher) handle(msg Msg) (err error) {
	switch msg.Code {
	case CodeSnapshotBlocks:
		bs := new(SnapshotBlocks)
		if err = bs.Deserialize(msg.Payload); err != nil {
			msg.Recycle()
			return err
		}
		msg.Recycle()

		f.log.Info(fmt.Sprintf("receive %d snapshotblocks %d from %s", len(bs.Blocks), msg.Id, msg.Sender))

		for _, block := range bs.Blocks {
			if err = f.receiver.receiveSnapshotBlock(block, types.RemoteFetch); err != nil {
				return err
			}
		}

		if len(bs.Blocks) > 0 {
			f.filter.done(msg.Id, msg)
		}

	case CodeAccountBlocks:
		bs := new(AccountBlocks)
		if err = bs.Deserialize(msg.Payload); err != nil {
			msg.Recycle()
			return err
		}
		msg.Recycle()

		f.log.Info(fmt.Sprintf("receive %d accountblocks from %s", len(bs.Blocks), msg.Sender))

		for _, block := range bs.Blocks {
			if err = f.receiver.receiveAccountBlock(block, types.RemoteFetch); err != nil {
				return err
			}
		}

		if len(bs.Blocks) > 0 {
			f.filter.done(msg.Id, msg)
		}

	case CodeException:
		f.filter.fail(msg.Id, msg.Sender.Id)
	}

	return nil
}

func (f *fetcher) FetchSnapshotBlocks(hash types.Hash, count uint64) {
	if _, ok := f.blackBlocks[hash]; ok {
		return
	}

	if !f.st.syncExited() {
		f.log.Debug("in syncing flow, cannot fetch")
		return
	}

	// been suppressed
	r, hold := f.filter.hold(hash)
	if hold {
		f.log.Debug(fmt.Sprintf("fetch suppressed GetSnapshotBlocks[hash %s, count %d]", hash, count))
		return
	}

	ps := f.filter.pickTargets(r, 0, f.peers)

	if len(ps) == 0 {
		f.log.Warn(fmt.Sprintf("no suit peers for %s", hash))
		return
	}

	for _, p := range ps {
		if p != nil {
			m := &GetSnapshotBlocks{
				From:    ledger.HashHeight{Hash: hash},
				Count:   count,
				Forward: false,
			}

			if err := p.send(CodeGetSnapshotBlocks, r.id, m); err != nil {
				f.log.Error(fmt.Sprintf("failed to send GetSnapshotBlocks[hash %s, count %d] to %s: %v", hash, count, p, err))
				f.filter.fail(r.id, p.Id)
			} else {
				f.log.Info(fmt.Sprintf("send GetSnapshotBlocks[hash %s, count %d] to %s", hash, count, p))
			}
		} else {
			f.log.Error(fmt.Sprintf("failed to fetch GetSnapshotBlocks[hash %s, count %d]: %v", hash, count, errNoSuitablePeer))
			f.filter.fail(r.id, peerId{})
		}
	}
}

// FetchSnapshotBlocksWithHeight fetch blocks:
//  ... count blocks ... {hash, height}
func (f *fetcher) FetchSnapshotBlocksWithHeight(hash types.Hash, height uint64, count uint64) {
	if _, ok := f.blackBlocks[hash]; ok {
		return
	}

	if !f.st.syncExited() {
		f.log.Debug("in syncing flow, cannot fetch")
		return
	}

	r, hold := f.filter.hold(hash)
	// been suppressed
	if hold {
		f.log.Debug(fmt.Sprintf("fetch suppressed GetSnapshotBlocks[hash %s, count %d]", hash, count))
		return
	}

	ps := f.filter.pickTargets(r, height, f.peers)

	if len(ps) == 0 {
		f.log.Warn(fmt.Sprintf("no suit peers for %s", hash))
		return
	}

	for _, p := range ps {
		if p != nil {
			m := &GetSnapshotBlocks{
				From:    ledger.HashHeight{Hash: hash},
				Count:   count,
				Forward: false,
			}

			if err := p.send(CodeGetSnapshotBlocks, r.id, m); err != nil {
				f.log.Error(fmt.Sprintf("failed to send GetSnapshotBlocks[hash %s, count %d] to %s: %v", hash, count, p, err))
				f.filter.fail(r.id, p.Id)
			} else {
				f.log.Info(fmt.Sprintf("send GetSnapshotBlocks[hash %s, count %d] to %s", hash, count, p))
			}
		} else {
			f.log.Error(fmt.Sprintf("failed to fetch GetSnapshotBlocks[hash %s, count %d]: %v", hash, count, errNoSuitablePeer))
			f.filter.fail(r.id, peerId{})
		}
	}
}

func (f *fetcher) FetchAccountBlocks(start types.Hash, count uint64, address *types.Address) {
	if _, ok := f.blackBlocks[start]; ok {
		return
	}

	if !f.st.syncExited() {
		f.log.Debug("in syncing flow, cannot fetch")
		return
	}

	r, hold := f.filter.hold(start)
	// been suppressed
	if hold {
		f.log.Debug(fmt.Sprintf("fetch suppressed GetAccountBlocks[hash %s, count %d]", start, count))
		return
	}

	ps := f.filter.pickTargets(r, 0, f.peers)

	if len(ps) == 0 {
		f.log.Warn(fmt.Sprintf("no suit peers for %s", start))
		return
	}

	for _, p := range ps {
		if p != nil {
			addr := ZERO_ADDRESS
			if address != nil {
				addr = *address
			}
			m := &GetAccountBlocks{
				Address: addr,
				From: ledger.HashHeight{
					Hash: start,
				},
				Count:   count,
				Forward: false,
			}

			if err := p.send(CodeGetAccountBlocks, r.id, m); err != nil {
				f.log.Error(fmt.Sprintf("failed to send GetAccountBlocks[hash %s, count %d] to %s: %v", start, count, p, err))
				f.filter.fail(r.id, p.Id)
			} else {
				f.log.Info(fmt.Sprintf("send GetAccountBlocks[hash %s, count %d] to %s", start, count, p))
			}
		} else {
			f.log.Error(fmt.Sprintf("failed to fetch GetAccountBlocks[hash %s, count %d]: %v", start, count, errNoSuitablePeer))
			f.filter.fail(r.id, peerId{})
		}
	}
}

func (f *fetcher) FetchAccountBlocksWithHeight(start types.Hash, count uint64, address *types.Address, sHeight uint64) {
	if _, ok := f.blackBlocks[start]; ok {
		return
	}

	if !f.st.syncExited() {
		f.log.Debug("in syncing flow, cannot fetch")
		return
	}

	r, hold := f.filter.hold(start)
	// been suppressed
	if hold {
		f.log.Debug(fmt.Sprintf("fetch suppressed GetAccountBlocks[hash %s, count %d]", start, count))
		return
	}

	ps := f.filter.pickTargets(r, sHeight, f.peers)

	if len(ps) == 0 {
		f.log.Warn(fmt.Sprintf("no suit peers for %s", start))
		return
	}

	for _, p := range ps {
		if p != nil {
			addr := ZERO_ADDRESS
			if address != nil {
				addr = *address
			}
			m := &GetAccountBlocks{
				Address: addr,
				From: ledger.HashHeight{
					Hash: start,
				},
				Count:   count,
				Forward: false,
			}

			if err := p.send(CodeGetAccountBlocks, r.id, m); err != nil {
				f.log.Error(fmt.Sprintf("failed to send GetAccountBlocks[hash %s, count %d] to %s: %v", start, count, p, err))
				f.filter.fail(r.id, p.Id)
			} else {
				f.log.Info(fmt.Sprintf("send GetAccountBlocks[hash %s, count %d] to %s", start, count, p))
			}
		} else {
			f.log.Error(fmt.Sprintf("failed to fetch GetAccountBlocks[hash %s, count %d]: %v", start, count, errNoSuitablePeer))
			f.filter.fail(r.id, peerId{})
		}
	}
}

func (f *fetcher) fetchSnapshotBlock(hash types.Hash, peer *Peer, callback func(msg Msg, err error)) {
	m := &GetSnapshotBlocks{
		From: ledger.HashHeight{
			Hash: hash,
		},
		Count:   1,
		Forward: false,
	}

	r := f.filter.add(hash)
	r.callback = callback

	err := peer.send(CodeGetSnapshotBlocks, r.id, m)
	if err != nil {
		r.done(Msg{}, err)
		f.log.Warn(fmt.Sprintf("failed to query reliable %s to %s", hash, peer))
	} else {
		f.log.Info(fmt.Sprintf("query reliable %s %d to %s", hash, r.id, peer))
	}
}
