package main

import "fmt"

func main(){
	n :=10
	rs:=Shuixianhua(n)
	fmt.Println(rs)
}
//f(n)=f(n-1)+f(n-2) {f(1)=1,f(2)=1+f(0)=1,f(0)=0}
func Shuixianhua(n int)int{
	if n<=1 {
		return n
	}

	return Shuixianhua(n-1)+Shuixianhua(n-2)

}