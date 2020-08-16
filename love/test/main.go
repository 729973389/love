package main

import (
	"fmt"
	"log"
)

func main() {


	GetFlag(1,100,100)
	
}

func GetFlag(head int,tail int,dup int){
	flag :=0
	min:=Max(head,tail)
	a:=0
	var gap int
	total:=(head+tail)*(tail-head+1)/2
	log.Println(total)
	for i:=head;i<=tail;i++ {
		a+=i
		gap=total-2*a
		if gap<0 {
			break
		}
		if gap<min {
			gap=min
			flag=i
		}
	}
	if flag != 0&&flag!=dup{
		fmt.Println(flag)
		GetFlag(head,flag,flag)
		GetFlag(flag,tail,flag)
	}



}

func Max(a,b int)int{
	if a>b{
		return a
	}
	return b

}