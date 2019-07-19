package main

import (
	"github.com/koinotice/vite"
	"github.com/koinotice/vite/cmd/gvite_plugins"
	_ "net/http/pprof"
)

// gvite is the official command-line client for Vite

func main() {
	govite.PrintBuildVersion()
	gvite_plugins.Loading()
}
