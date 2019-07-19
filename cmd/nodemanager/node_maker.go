package nodemanager

import (
	"github.com/koinotice/vite/node"
	"gopkg.in/urfave/cli.v1"
)

type NodeMaker interface {

	//create Node
	MakeNode(ctx *cli.Context) (*node.Node, error)

	//create NodeConfig
	MakeNodeConfig(ctx *cli.Context) (*node.Config, error)
}
