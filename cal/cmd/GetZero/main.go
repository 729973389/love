package main

import "fmt"

func main() {
	fmt.Println(GetZero(100000))

}

func GetZero(n int) int {
	count10 := 0
	count2, count5, min := 0, 0, 0
	for i := 1; i <= n; i++ {
		j := i
		for j%10 == 0 {
			j /= 10
			count10++
		}
		for j%2 == 0 {
			j /= 2
			count2++
		}
		for j%5 == 0 {
			j /= 5
			count5++
		}
	}
	min = count2
	if count2 > count5 {
		min = count5
	}
	return count10 + min
}
