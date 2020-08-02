package main

import (
	"flag"
	"fmt"
	"strings"
)

var s = flag.String("s"," ","echo what you write ")

func main(){
	flag.Parse()
	fmt.Println(strings.Join(flag.Args(),*s))
	for _,v := range flag.Args() {
		fmt.Println(v)
	}
}

