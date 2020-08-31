package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {
	var b [1024]byte
	defer func() {
		n := runtime.Stack(b[:], true)
		log.Println(string(b[:n]))
		log.Println("It's so beautiful!!!!")
	}()
	defer log.Println("I'm last!")
	defer MyPrint()()
	fmt.Println(Add(1, 2, 3, 4, 5, 6))
	s := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(Add(s...))
	fmt.Printf("%T", Add)
	D := Double(4)
	fmt.Println(D)
	trible := Trible
	fmt.Println(trible(4))
	log.Println(Return(100))
}
func Return(x int) (result int) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		result = x
	}()
	panic("f")
}

func MyPrint() func() {
	log.Println("START")
	return func() {
		time.Sleep(100 * time.Millisecond)
		log.Println("EXIT")
	}
}
func Double(x int) (result int) {
	defer func() { fmt.Println(result) }()
	return x + x
}
func Trible(x int) (result int) {
	defer func() { result += x }()
	return x + x
}
func Add(num ...int) int {
	var total int
	for _, v := range num {
		total += v
	}
	return total
}

//func crawl(url string) []string {
//	resp, err := http.Get(url)
//	if err != nil {
//		log.Error(errors.Wrap(err, "GetHttp"))
//		return nil
//	}
//	if resp.StatusCode != http.StatusOK {
//		resp.Body.Close()
//		return nil
//	}
//	doc, err := html.Parse(resp.Body)
//	if err != nil {
//		return nil
//	}
//	resp.Body.Close()
//	var links []string
//	visitNode := func(n *html.Node) {
//		if n.Type == html.ElementNode && n.Data == "a" {
//			for _, a := range n.Attr {
//				if a.Key != "href" {
//					continue
//				}
//				link, err := resp.Request.URL.Parse(a.Val)
//				if err != nil {
//					continue
//				}
//				links = append(links, link.String())
//			}
//		}
//
//	}
//
//}

