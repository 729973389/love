package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(Reserve(a))
	Reserve2(a, 4)
	Reserve2(a, len(a))
	fmt.Println(a)

}

func Reserve(n []int) []int {
	b := make([]int, 0)
	for i := len(n) - 1; i >= 0; i-- {
		b = append(b, n[i])
	}
	return b

}

func Reserve2(n []int, length int) {
	for i := 0; i < length; i++ {
		if i > (length-1)/2 {
			break
		}
		n[i], n[length-1-i] = n[length-1-i], n[i]
	}

}
