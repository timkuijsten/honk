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
	"golang.org/x/sys/unix"
)

func xunveil(path string, perms string) {
	err := unix.Unveil(path, perms)
	if err != nil {
		panic(fmt.Errorf("unveil(%s, %s) failure (%d)", path, perms, err))
	}
}

func xpledge(promises string) {
	err := unix.PledgePromises(promises)
	if err != nil {
		panic(fmt.Errorf("pledge(%s) failure (%d)", promises, err))
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
