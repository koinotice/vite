package vite

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/koinotice/vite/wallet/hd-bip/derivation"

	"github.com/koinotice/vite/chain"
	"github.com/koinotice/vite/common/fork"
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/config"
	"github.com/koinotice/vite/consensus"
	"github.com/koinotice/vite/log15"
	"github.com/koinotice/vite/net"
	"github.com/koinotice/vite/onroad"
	"github.com/koinotice/vite/pool"
	"github.com/koinotice/vite/producer"
	"github.com/koinotice/vite/verifier"
	"github.com/koinotice/vite/vm"
	"github.com/koinotice/vite/wallet"
)

var (
	log = log15.New("module", "console/bridge")
)

type Vite struct {
	config *config.Config

	walletManager   *wallet.Manager
	accountVerifier verifier.Verifier
	chain           chain.Chain
	producer        producer.Producer
	net             net.Net
	pool            pool.BlockPool
	consensus       consensus.Consensus
	onRoad          *onroad.Manager
}

func New(cfg *config.Config, walletManager *wallet.Manager) (vite *Vite, err error) {
	var addressContext *producer.AddressContext
	if cfg.Producer.Producer && cfg.Producer.Coinbase != "" {
		var coinbase *types.Address
		var index uint32
		coinbase, index, err = parseCoinbase(cfg.Producer.Coinbase)
		if err != nil {
			log.Error(fmt.Sprintf("coinBase parse fail. %v", cfg.Producer.Coinbase), "err", err)
			return nil, err
		}
		err = walletManager.MatchAddress(cfg.EntropyStorePath, *coinbase, index)

		if err != nil {
			log.Error(fmt.Sprintf("coinBase is not child of entropyStore, coinBase is : %v", cfg.Producer.Coinbase), "err", err)
			return nil, err
		}

		var key *derivation.Key
		_, key, _, err = walletManager.GlobalFindAddr(*coinbase)
		if err != nil {
			return
		}

		cfg.Net.MineKey, err = key.PrivateKey()
		if err != nil {
			return
		}

		addressContext = &producer.AddressContext{
			EntryPath: cfg.EntropyStorePath,
			Address:   *coinbase,
			Index:     index,
		}
	}

	// set fork points
	fork.SetForkPoints(cfg.ForkPoints)

	// chain
	chain := chain.NewChain(cfg.DataDir, cfg.Chain, cfg.Genesis)

	err = chain.Init()
	if err != nil {
		return nil, err
	}
	// pool
	pl, err := pool.NewPool(chain)
	if err != nil {
		return nil, err
	}
	// consensus
	cs := consensus.NewConsensus(chain, pl)

	// sb verifier
	aVerifier := verifier.NewAccountVerifier(chain, cs)
	sbVerifier := verifier.NewSnapshotVerifier(chain, cs)

	verifier := verifier.NewVerifier(sbVerifier, aVerifier)
	// net
	net, err := net.New(cfg.Net, chain, verifier, cs, pl)
	if err != nil {
		return
	}

	// vite
	vite = &Vite{
		config:          cfg,
		walletManager:   walletManager,
		net:             net,
		chain:           chain,
		pool:            pl,
		consensus:       cs,
		accountVerifier: verifier,
	}

	if addressContext != nil {
		vite.producer = producer.NewProducer(chain, net, addressContext, cs, sbVerifier, walletManager, pl)
	}

	// onroad
	or := onroad.NewManager(net, pl, vite.producer, vite.consensus, walletManager)

	// set onroad
	vite.onRoad = or
	return
}

func (v *Vite) Init() (err error) {
	vm.InitVMConfig(v.config.IsVmTest, v.config.IsUseVmTestParam, v.config.IsVmDebug, v.config.DataDir)

	//v.chain.Init()
	if v.producer != nil {
		if err := v.producer.Init(); err != nil {
			log.Error("Init producer failed, error is "+err.Error(), "method", "vite.Init")
			return err
		}
	}

	// initOnRoadPool
	v.onRoad.Init(v.chain)
	v.accountVerifier.InitOnRoadPool(v.onRoad)

	return nil
}

func (v *Vite) Start() (err error) {
	v.onRoad.Start()

	v.chain.Start()

	err = v.consensus.Init()
	if err != nil {
		return err
	}
	// hack
	v.pool.Init(v.net, v.walletManager, v.accountVerifier.GetSnapshotVerifier(), v.accountVerifier, v.consensus)

	v.consensus.Start()

	err = v.net.Start()
	if err != nil {
		return
	}

	v.pool.Start()
	if v.producer != nil {

		if err := v.producer.Start(); err != nil {
			log.Error("producer.Start failed, error is "+err.Error(), "method", "vite.Start")
			return err
		}
	}
	return nil
}

func (v *Vite) Stop() (err error) {

	v.net.Stop()
	v.pool.Stop()

	if v.producer != nil {
		if err := v.producer.Stop(); err != nil {
			log.Error("producer.Stop failed, error is "+err.Error(), "method", "vite.Stop")
			return err
		}
	}
	v.consensus.Stop()
	v.chain.Stop()
	v.onRoad.Stop()
	return nil
}

func (v *Vite) Chain() chain.Chain {
	return v.chain
}

func (v *Vite) Net() net.Net {
	return v.net
}

func (v *Vite) WalletManager() *wallet.Manager {
	return v.walletManager
}

func (v *Vite) Producer() producer.Producer {
	return v.producer
}

func (v *Vite) Pool() pool.BlockPool {
	return v.pool
}

func (v *Vite) Consensus() consensus.Consensus {
	return v.consensus
}

func (v *Vite) OnRoad() *onroad.Manager {
	return v.onRoad
}

func (v *Vite) Config() *config.Config {
	return v.config
}

func parseCoinbase(coinbaseCfg string) (*types.Address, uint32, error) {
	splits := strings.Split(coinbaseCfg, ":")
	if len(splits) != 2 {
		return nil, 0, errors.New("len is not equals 2.")
	}
	i, err := strconv.Atoi(splits[0])
	if err != nil {
		return nil, 0, err
	}
	addr, err := types.HexToAddress(splits[1])
	if err != nil {
		return nil, 0, err
	}

	return &addr, uint32(i), nil
}
