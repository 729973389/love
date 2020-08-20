// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//  ignore

package main

import (
	"github.com/wuff1996/edgeHub/cmd/root"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go root.Run()
	wg.Wait()
}
