package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	ParameterStatement()  //goroutine:  out \n goroutine:  inner \n function:  out
	ParameterLifecycle2() //EXIT: ParameterLifecycle \n 1 2 3 4 5 \n
}

//statement the same name in goroutine,the parameter will cover outer parameter
func ParameterStatement() {
	parameter := "out"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("goroutine: ", parameter)
		parameter := "inner"
		fmt.Println("goroutine: ", parameter)
	}()
	wg.Wait()
	fmt.Println("function: ", parameter)
	fmt.Printf("EXIT: ParameterStatement\n\n")
}

//when go routine context func exit the parameter's value is saved inside the goroutine and it will not exit except the main
//process exit ,because when the main process exit ,the program will end running.
func ParameterLifecycle(wg *sync.WaitGroup) {
	defer fmt.Printf("EXIT: ParameterLifecycle\n\n")
	a := 0
	timer := time.NewTimer(5 * time.Second)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			select {
			case <-timer.C:
				a++
				log.Println(a)
				timer.Reset(5 * time.Second)
			}
		}
	}()
}

//help ParameterLifecycle
func ParameterLifecycle2() {
	defer fmt.Printf("EXIT: ParameterLifecycle2\n\n")
	var wg sync.WaitGroup
	wg.Add(1)
	go ParameterLifecycle(&wg)
	wg.Wait()
}
