package main

import (
	"log"
	"sync"
	"time"
)

//type you string
func main() {
	var ch =make(chan []byte,256)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for  {
			ch<-nil
			time.Sleep(5 * time.Second)

		}
	}()

	go func() {
		for  {
			select {
			case b,ok:=<-ch:
				if !ok {
					log.Println("false")
				}
				if b==nil{
					log.Println("nil")
				}
				log.Println("receive",b)
			}

		}
	}()
	wg.Wait()
}

//func Love(you you)