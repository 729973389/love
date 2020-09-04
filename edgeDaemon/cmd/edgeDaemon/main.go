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
	log.Info("Version: ", Version)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	Mysignal := make(chan os.Signal, 1)
	signal.Notify(Mysignal)
	go func() {
		for {
			select {
			case s := <-Mysignal:
				go func() {
					switch s {
					case os.Interrupt:
						log.Info("Interrupt")
						cancel()
					case os.Kill:
						log.Info("Kill")
						cancel()
					default:
						cancel()
					}
				}()
			case <-ctx.Done():
				ticker := time.NewTicker(10 * time.Second)
				select {
				case <-ticker.C:
					log.Warning("EXIT : MAIN : ABNORMAL")
					os.Exit(1)
				}
			}

		}
	}()
	hub := root.NewHub()
	wg.Add(1)
	go func() {
		defer wg.Done()
		root.RunWS(ctx, hub)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		root.RunMQTT(ctx, hub)
	}()
	wg.Wait()
	cancel()
	log.Println("MAIN : EXIT: NORMAL")
}
