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
				go func() {
					switch s {
					case os.Kill:
						log.Info("KILL")
						cancel()
					case os.Interrupt:
						log.Info("INTERRUPT")
						cancel()
					}
				}()
			case <-ctx.Done():
				ticker := time.NewTicker(10 * time.Second)
				select {
				case <-ticker.C:
					log.Error("EXIT: MAIN: ABNORMAL")
					os.Exit(1)
				}
			}
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		root.Run(ctx)
	}()
	wg.Wait()
	log.Info("EXIT: MAIN: NORMAL")
}
