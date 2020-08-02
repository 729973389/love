package main

import (
	"fmt"
	"strings"
)

func main () {
	str := "31414-4324:3243 324"
	d :=str[:0]
	for _,v:=range str{
		if strings.ContainsAny(string(v),"123456789") {
			d += string(v)
		}
		fmt.Println(v)
	}
	fmt.Println(d)
}