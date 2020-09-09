package main

import "fmt"

func main() {
	ForContinue() //continue
	ForBreak()    //break  \n EXIT: ForBreak
}

//this acts as the last forCycle,it shows that the continue at last time will act like return,it directly exits the
//function body
func ForContinue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %v", p)
		}
	}()
	for i := 0; i < 1; i++ {
		fmt.Println("continue")
		continue
	}
	fmt.Sprintf("EXIT: ForContinue\n\n")
}

//this acts as the last forCycle,shows that the break at last just break the forCycle and do the task after forCycle
func ForBreak() {
	for i := 0; i < 1; i++ {
		fmt.Println("break")
		break
	}
	fmt.Printf("EXIT: ForBreak\n\n")
}

