package main

import "fmt"

const k = 2

func main() {
	getK := GetKestNumber()
	fmt.Println(getK(1))
	fmt.Println(getK(2))
	fmt.Println(getK(0))
	fmt.Println(getK(5))
	fmt.Println(getK(6))
	fmt.Println(getK(0))
	fmt.Println(getK(1))
	fmt.Println(getK(2))
	fmt.Println(getK(0))

}

func GetKestNumber() func(n int) int {
	var slice = make([]int, 0)
	return func(n int) int {
		if len(slice) < k+1 {
			slice = append(slice, n)
		} else {
			slice[k] = n
		}
		for i := len(slice) - 1; i > 0; i-- {
			if slice[i] > slice[i-1] {
				slice[i], slice[i-1] = slice[i-1], slice[i]
			}
		}
		if len(slice) < k {
			return slice[len(slice)-1]
		}
		return slice[k-1]
	}

}
