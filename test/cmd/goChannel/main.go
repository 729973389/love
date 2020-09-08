package main

import (
	"fmt"
	"time"
)

func main() {
	ChannelRace()         //*************  or ---------------
	TimerChannel()        //time out
	CloseReadChannel()    //read
	CloseWriteChannel()   //panic: send on closed channel
	DefaultChannel()      //default
	SelectBreakContinue() //break     I'm after select :)  \n continue            /

}

//race between ch <- GetInt() and <-timer.C.The result shows that the ch <- GetInt() is not valid ,th select
//occasionally chose one of them.
func ChannelRace() {
	var ch = make(chan int, 1)
	timer := time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("*************")
	case ch <- GetInt():
		fmt.Println("---------------")
	}
	fmt.Printf("EXIT: ChannelRace\n\n")
}

// func to get int number
func GetInt() (n int) {
	n = 66
	time.Sleep(3 * time.Second)
	return
}

//to do a task in a limit time
func TimerChannel() {
	const timeLimit = 2 * time.Second
	ticker := time.NewTicker(timeLimit)
	defer ticker.Stop()
	var sigch = make(chan *struct{})
	var n int
	go func() {
		n = GetInt()
		sigch <- &struct{}{}
	}()
	select {
	case <-ticker.C:
		fmt.Println("time out")
	case <-sigch:
		fmt.Println(n)
	}
	fmt.Printf("EXIT: TimerChannel\n\n")
}

//close channel when some task is reading
func CloseReadChannel() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()
	var ch = make(chan string)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	go func() {
		select {
		case <-ticker.C:
			close(ch)
		}
	}()
	select {
	case <-ch:
		fmt.Println("read")
	}
	fmt.Printf("EXIT: CloseReadChannel\n\n")
}

//close channel when some task is writing to it
func CloseWriteChannel() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %v\n", p)
		}
	}()
	var ch = make(chan string)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	go func() {
		select {
		case <-ticker.C:
			close(ch)
		}
	}()
	select {
	case ch <- "write":
		fmt.Println("write")
	}
	fmt.Printf("EXIT: CloseWriteChannel\n\n")
}

//default body is always chosen in context of select
func DefaultChannel() {
	var ch = make(chan *struct{})
	select {
	case ch <- &struct{}{}:
		fmt.Println("chan")
	default:
		fmt.Println("default")
	}
	fmt.Printf("EXIT DefaultChannel\n\n")
}

//break & continue in select,so it shows that break can just break out the select,it can't break out for cycle,what's more
//continue can ignore what behind and continue the next for cycle.
func SelectBreakContinue() {
	selectBreak := make(chan string)
	selectContinue := make(chan string)
	go func() {
		selectBreak <- "break"
		selectContinue <- "continue"
	}()
	for i := 0; i < 2; i++ {
		select {
		case s := <-selectContinue:
			fmt.Printf("%-10s\n", s)
			continue
		case s := <-selectBreak:
			fmt.Printf("%-10s", s)
			break
		}
		fmt.Println("I'm after select :)")
	}

	fmt.Printf("EXIT: SelectBreakContinue\n\n")
}
