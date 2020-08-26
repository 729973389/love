package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeDaemon/cmd/edgeDaemon/root"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	Version = "v1"
	Build   = "N/A"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	Mysignal := make(chan os.Signal, 1)
	signal.Notify(Mysignal)
	go func() {
		for {
			select {
			case s := <-Mysignal:
				switch s {
				case os.Interrupt:
					cancel()
				case os.Kill:
					cancel()
				}

			}

		}
	}()
	wg.Add(1)
	hub := root.NewHub()
	go root.RunTCP(ctx, &wg, hub)
	wg.Wait()
	cancel()
	time.Sleep(10 * time.Second)
	log.Println("EXIT")
}
