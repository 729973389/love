package main

import (
	"fmt"
)

func main() {
	a := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(BinarySort(a, 0, len(a)-1, 9))
}

func BinarySort(n [10]int, head, tail, target int) int {
	if head > tail {
		return -1
	}
	//如果使用(tail+head)/2,当tail和head足够大时，int类型就越界了，所以得(tail-head)/2+head
	mid := (tail-head)/2 + head
	fmt.Println(n[mid])
	if n[mid] == target {
		return mid
	} else if n[mid] < target {
		return BinarySort(n, mid+1, tail, target)
	} else if n[mid] > target {
		return BinarySort(n, head, mid-1, target)
	}
	return -1
}
