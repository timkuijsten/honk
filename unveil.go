//go:build openbsd
// +build openbsd

//
// Copyright (c) 2019 Ted Unangst <tedu@tedunangst.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	SYS_PLEDGE = 108
	SYS_UNVEIL = 114
)

func xunveil(path, perms string) {
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		panic(fmt.Errorf("unveil(%s, %s) failure (%v)", path, perms, err))
	}
	permsPtr, err := syscall.BytePtrFromString(perms)
	if err != nil {
		panic(fmt.Errorf("unveil(%s, %s) failure (%v)", path, perms, err))
	}

	_, _, e := syscall.Syscall(SYS_UNVEIL, uintptr(unsafe.Pointer(pathPtr)), uintptr(unsafe.Pointer(permsPtr)), 0)
	if e != 0 {
		panic(fmt.Errorf("unveil(%s, %s) failure (%d)", path, perms, e))
	}
}

func xpledge(promises string) {
	var expr unsafe.Pointer

	pptr, err := syscall.BytePtrFromString(promises)
	if err != nil {
		panic(fmt.Errorf("pledge(%s) failure (%v)", promises, err))
	}

	_, _, e := syscall.Syscall(SYS_PLEDGE, uintptr(unsafe.Pointer(pptr)), uintptr(expr), 0)
	if e != 0 {
		panic(fmt.Errorf("pledge(%s) failure (%d)", promises, e))
	}
}

func init() {
	preservehooks = append(preservehooks, func() {
		xunveil("/etc/ssl", "r")
		if viewDir != dataDir {
			xunveil(viewDir, "r")
		}
		xunveil(dataDir, "rwc")
		xpledge("stdio rpath wpath cpath flock dns inet unix")
	})
	backendhooks = append(backendhooks, func() {
		xpledge("stdio unix")
	})
}
