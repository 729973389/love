package main

import (
	"github.com/wuff1996/edgeDaemon/cmd/root"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	root.RunTCP()
	wg.Wait()
}
