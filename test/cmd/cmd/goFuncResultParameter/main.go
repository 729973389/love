package main
import "fmt"

func main(){
	fmt.Println(ResultParameter("a","b"))
	fmt.Printf("%q",Test2())
}

func ResultParameter(a,b string)(s string){
	s=a+b
	return s
}

func Test2()(s string){

return s
}

