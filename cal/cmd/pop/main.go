package main

import (
	"fmt"
)

func main(){
	n := []int{1,5,6,7,8,78,45363,4234,65,4,67,78}
	Pop(n)
	fmt.Println(n)
}

func Pop(n []int){
	for i:=0;i<len(n);i++ {
		for j:=0;j<len(n)-1;j++ {
			if n[j]>n[j+1] {
				n[j],n[j+1]=n[j+1],n[j]
			}
		}

	}

}