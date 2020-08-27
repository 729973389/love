package main

import (
	"fmt"
)

func main(){
	n := []int{5,5,535,6346,53,1,5,67,7,8,3,4,33,6}
	ReplaceSort(n)
	fmt.Println(n)
}

func ReplaceSort(n []int){
	for i:=0;i<len(n);i++{
		minIndex:=i
		for j:=i;j<len(n);j++ {
			if n[j]<n[minIndex] {
				minIndex=j
			}
		}
		n[i],n[minIndex]=n[minIndex],n[i]
	}

}