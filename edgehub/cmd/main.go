// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//  ignore

package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/cmd/root"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	Version = "v1"
	Build   = "N/A"
)

const CloseTime = 10

func main() {
	log.Printf("Version: %s  started", Version)
	ctx, cancel := context.WithCancel(context.Background())
	Mysignal := make(chan os.Signal, 1)
	signal.Notify(Mysignal)
	go func() {
		for {
			select {
			case s := <-Mysignal:
				switch s {
				case os.Kill:
					cancel()
				case os.Interrupt:
					cancel()
				}
			}
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go root.Run(ctx, &wg)
	wg.Wait()
	cancel()
	time.Sleep(CloseTime * time.Second)
	log.Warning("Exit main")
}
