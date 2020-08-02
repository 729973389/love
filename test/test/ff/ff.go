package ff

import "fmt"

var A int32=a()

func a()int32 {
	fmt.Println("calling a()")
	return 1
}
func init(){
	fmt.Println("init ff")
}

