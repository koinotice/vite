package params

import "github.com/koinotice/vite"

var Version = func() string {
	return govite.VITE_BUILD_VERSION
}()
