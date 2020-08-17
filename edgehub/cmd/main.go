// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//  ignore

package main

import (
	tohttp "github.com/wuff1996/edgeHub/cmd/http"
	"github.com/wuff1996/edgeHub/cmd/root"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go root.Run()
	tohttp.GetConfig()
	s := tohttp.NewSchedule()
	go s.Run()
	go tohttp.Serve(s)
	wg.Wait()
}
