package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/koinotice/vite/net/vnode"
)

type mockSocket struct {
}

func (m *mockSocket) ping(n *Node, ch chan<- *Node) (err error) {
	go func() {
		n.ID = vnode.RandomNodeID()
		ch <- n
	}()

	return nil
}

func (m *mockSocket) pong(echo []byte, n *Node) (err error) {
	panic("implement me")
}

func (m *mockSocket) findNode(target vnode.NodeID, count int, n *Node, ch chan<- []*vnode.EndPoint) (err error) {
	go func() {
		var eps = make([]*vnode.EndPoint, count)
		for i := 0; i < count; i++ {
			eps[i] = &vnode.EndPoint{
				Host: []byte{0, 0, 0, 0},
				Port: int(i),
				Typ:  vnode.HostIPv4,
			}
		}

		ch <- eps
	}()

	return nil
}

func (m *mockSocket) sendNodes(eps []*vnode.EndPoint, addr *net.UDPAddr) (err error) {
	panic("implement me")
}

func (m *mockSocket) start() error {
	return nil
}

func (m *mockSocket) stop() error {
	return nil
}

func TestFindNode(t *testing.T) {
	tab := newTable(vnode.ZERO, self.Net, newListBucket, nil)
	tab.add(&Node{
		Node: vnode.Node{
			ID: vnode.RandFromDistance(tab.id, 100),
			EndPoint: vnode.EndPoint{
				Host: []byte{127, 0, 0, 1},
				Port: 8888,
				Typ:  vnode.HostIPv4,
			},
			Net: self.Net,
			Ext: nil,
		},
	})

	var d = &Discovery{
		node: &vnode.Node{
			ID: vnode.ZERO,
			EndPoint: vnode.EndPoint{
				Host: []byte{127, 0, 0, 1},
				Port: 8483,
				Typ:  vnode.HostIPv4,
			},
			Net: self.Net,
			Ext: nil,
		},
		table:  tab,
		finder: nil,
		stage:  make(map[string]*Node),
		socket: &mockSocket{},
	}

	nodes := d.lookup(vnode.ZERO, 32)
	if len(nodes) != tab.size() {
		t.Errorf("should not find %d nodes", len(nodes))
	}
}

// timer reset
func TestTimer(t *testing.T) {
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	<-timer.C

	if !timer.Stop() {
		//<-timer.C // will block
	}

	start := time.Now().Unix()
	timer.Reset(time.Second)
	<-timer.C

	if time.Now().Unix()-start != 1 {
		t.Fail()
	}
}
