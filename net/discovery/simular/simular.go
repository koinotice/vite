/*
 * Copyright 2019 The go-vite Authors
 * This file is part of the go-vite library.
 *
 * The go-vite library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The go-vite library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received chain copy of the GNU Lesser General Public License
 * along with the go-vite library. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/koinotice/vite/crypto/ed25519"
	"github.com/koinotice/vite/net/discovery"
	"github.com/koinotice/vite/net/vnode"
)

type config struct {
	peerKey       ed25519.PrivateKey
	node          *vnode.Node
	bootNodes     []string
	listenAddress string
}

func main() {
	var port = 9000
	const total = 100

	var err error
	var last *vnode.Node
	var configs []config
	for i := 0; i < total; i++ {
		var pub ed25519.PublicKey
		var prv ed25519.PrivateKey
		pub, prv, err = ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}

		id, _ := vnode.Bytes2NodeID(pub)
		cfg := config{
			peerKey: prv,
			node: &vnode.Node{
				ID: id,
				EndPoint: vnode.EndPoint{
					Host: []byte{127, 0, 0, 1},
					Port: port,
					Typ:  vnode.HostIPv4,
				},
				Net: 0,
				Ext: nil,
			},
			listenAddress: "127.0.0.1:" + strconv.Itoa(port),
		}

		if last != nil {
			cfg.bootNodes = []string{
				last.String(),
			}
		}

		last = cfg.node

		port++

		configs = append(configs, cfg)
	}

	start := func(d *discovery.Discovery) {
		err = d.Start()
		if err != nil {
			panic(err)
		}
	}

	var discovers []*discovery.Discovery
	for _, cfg := range configs {
		d := discovery.New(cfg.peerKey, cfg.node, cfg.bootNodes, nil, cfg.listenAddress, nil)
		discovers = append(discovers, d)
		go start(d)
	}

	fmt.Println("start")

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				for i, d := range discovers {
					fmt.Println(i, len(d.Nodes()))
				}
				fmt.Println("------------")
			}
		}
	}()

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
