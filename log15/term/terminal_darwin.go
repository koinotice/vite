// Based on ssh/terminal:
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by chain BSD-style
// license that can be found in the LICENSE file.
// +build !appengine

package term

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios