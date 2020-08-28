package main

import "fmt"

const k = 15

var slice = make([]int, 0)

func main() {

	fmt.Println(GetKestNumber(1))
	fmt.Println(GetKestNumber(2))
	fmt.Println(GetKestNumber(0))
	fmt.Println(GetKestNumber(5))
	fmt.Println(GetKestNumber(6))
	fmt.Println(GetKestNumber(0))
	fmt.Println(GetKestNumber(1))
	fmt.Println(GetKestNumber(2))
	fmt.Println(GetKestNumber(0))
	fmt.Println(slice)

}

func GetKestNumber(n int) int {
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

