package discovery

import (
	"net"
	"time"

	"github.com/koinotice/vite/net/netool"
	"github.com/koinotice/vite/net/vnode"
)

type Node struct {
	vnode.Node
	checkAt  int64
	addAt    int64
	activeAt int64
	checking bool // is in check flow
	finding  bool // is finding some target from this node
	addr     *net.UDPAddr
	parseAt  int64 // last time addr parsed
}

func (n *Node) udpAddr() (addr *net.UDPAddr, err error) {
	now := time.Now().Unix()

	// 15min
	if now-n.parseAt > 900 || n.addr == nil {
		addr, err = net.ResolveUDPAddr("udp", n.Address())
		if err != nil {
			return
		}

		n.addr = addr
		n.parseAt = now
		return
	}

	return n.addr, nil
}

func (n *Node) update(n2 *Node) {
	n.ID = n2.ID
	n.Ext = n2.Ext
	n.Net = n2.Net
	n.EndPoint = n2.EndPoint
}

func extractEndPoint(addr *net.UDPAddr, from *vnode.EndPoint) (e *vnode.EndPoint, addr2 *net.UDPAddr) {
	var err error
	var done bool
	if from != nil {
		// from is available
		addr2, err = net.ResolveUDPAddr("udp", from.String())
		if err == nil {
			if from.Typ.Is(vnode.HostDomain) || netool.CheckRelayIP(addr.IP, from.Host) == nil {
				// from is domain, or IP is available
				done = true
				e = from
			}
		}
	}

	if !done {
		e = udpAddrToEndPoint(addr)
		addr2 = addr
	}

	return
}

func nodeFromEndPoint(e vnode.EndPoint) (n *Node, err error) {
	udp, err := net.ResolveUDPAddr("udp", e.String())
	if err != nil {
		return
	}

	return &Node{
		Node: vnode.Node{
			EndPoint: e,
		},
		addr:    udp,
		parseAt: time.Now().Unix(),
	}, nil
}

func udpAddrToEndPoint(addr *net.UDPAddr) (e *vnode.EndPoint) {
	e = new(vnode.EndPoint)
	if ip4 := addr.IP.To4(); len(ip4) != 0 {
		e.Host = ip4
		e.Typ = vnode.HostIPv4
	} else {
		e.Host = addr.IP
		e.Typ = vnode.HostIPv6
	}
	e.Port = addr.Port

	return
}

func nodeFromPing(res *packet) *Node {
	p := res.body.(*ping)

	e, addr := extractEndPoint(res.from, p.from)

	return &Node{
		Node: vnode.Node{
			ID:       res.id,
			EndPoint: *e,
			Net:      p.net,
			Ext:      p.ext,
		},
		addr:    addr,
		parseAt: time.Now().Unix(),
	}
}

func nodeFromPong(res *packet) *Node {
	p := res.body.(*pong)

	e, addr := extractEndPoint(res.from, p.from)

	return &Node{
		Node: vnode.Node{
			ID:       res.id,
			EndPoint: *e,
			Net:      p.net,
			Ext:      p.ext,
		},
		addr:    addr,
		parseAt: time.Now().Unix(),
	}
}
