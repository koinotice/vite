package nodemanager

import "github.com/koinotice/vite/node"

type NodeManager interface {
	Start() error

	Stop() error

	Node() *node.Node
}
