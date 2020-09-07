package main

import (
		"context"
			"fmt"
				"github.com/pkg/errors"
					"log"
						"strings"
							"sync"
								"time"
							)

							func main() {
									ctx,cancel:=context.WithCancel(context.Background())
										//var l=Loop1(&wg)
											var wg sync.WaitGroup
												wg.Add(1)
													//wg.Add(1)
														//go func() {
															//	wg.Done()
																//	l()
																	//	fmt.Println("exit goroutine")
																		//}()
																			var l2=loop2(ctx)
																				go func() {

																							l2()
																									cancel()
																										}()
																											wg.Wait()
																												fmt.Println("exit")
																												    //ch := make(chan *struct{})
																												        //go func() {
																														//	for  {
																															//		ch<-&struct{}{}
																																//		time.Sleep(5 * time.Second)
																																	//
																																		//	}
																																			//}()
																																				//
																																				    //for {
																																				    	//	select {
																																						//	case a:=<-ch:
																																							//		fmt.Println(a)
																																								//
																																									//	}
																																										//}
																																											//c:=1
																																												//var wg sync.WaitGroup
																																													//wg.Add(1)
																																														//go func() {
																																															//	defer wg.Done()
																																																//	c=2
																																																	//	c:=3
																																																		//	fmt.Println("goroutine: ",c)
																																																			//}()
																																																				//wg.Wait()
																																																					//fmt.Println(c)

																																																						//
																																																							//timer :=time.NewTimer(10 * time.Second)
																																																								//defer timer.Stop()
																																																									//timer.Stop()
																																																										//s:="{\"id\":\"lxd\",\"token\":\"hello\"}"
																																																											//fmt.Println(FindKeyString(s,"d"))
																																																												//fmt.Println(goJSON.FindKeyString(s,"d"))
																																																													//var wg sync.WaitGroup
																																																														//	var ch = make (chan int)
																																																															//	wg.Add(1)
																																																																//	go func() {
																																																																	//		defer wg.Done()
																																																																		//		for  {
																																																																			//			ch<-1
																																																																				//			time.Sleep(10 * time.Minute)
																																																																					//		}
																																																																						//	}()
																																																																							//	for  {
																																																																								//		select {
																																																																									//		case <-ch:
																																																																										//			log.Println("Receive")
																																																																											//		}
																																																																												//		wg.Wait()
																																																																													//		fmt.Println("EXIT")
																																																																														//
																																																																															//	}
																																																																																//defer func() {
																																																																																	//	fmt.Println("A:)")
																																																																																		//	fmt.Println("B:(")
																																																																																			//}()
																																																																																				//a := make(chan string)
																																																																																					//b := make(chan string)
																																																																																						//go func() {
																																																																																							//	for {
																																																																																								//		a <- "break"
																																																																																									//		time.Sleep(10 * time.Second)
																																																																																										//		b <- "continue"
																																																																																											//	}
																																																																																												//}()
																																																																																													//for {
																																																																																														//	select {
																																																																																															//	case s := <-a:
																																																																																																//		fmt.Printf("%-10s\n", s)
																																																																																																	//		continue
																																																																																																		//	case s := <-b:
																																																																																																			//		fmt.Printf("%-10s", s)
																																																																																																				//		break
																																																																																																					//	}
																																																																																																						//	fmt.Println(":)")
																																																																																																							//}
																																																																																																						}

																																																																																																						func p1() {
																																																																																																								defer fmt.Println("EXIT p1")
																																																																																																									a := 0
																																																																																																										timer := time.NewTimer(5 * time.Second)
																																																																																																											fmt.Println("1")
																																																																																																												go func() {
																																																																																																															for {
																																																																																																																			select {
																																																																																																																							case <-timer.C:
																																																																																																																												a++
																																																																																																																																log.Println(a)
																																																																																																																																				timer.Reset(5 * time.Second)
																																																																																																																																							}
																																																																																																																																									}
																																																																																																																																										}()
																																																																																																																																									}

																																																																																																																																									func p2() {
																																																																																																																																											defer fmt.Println("EXIT p2")
																																																																																																																																												fmt.Println("2")
																																																																																																																																													go p1()
																																																																																																																																												}
																																																																																																																																												func FindKeyString(s string, key string) (string, error) {
																																																																																																																																														if !strings.Contains(s, ",") {
																																																																																																																																																	if strings.Contains(s, "\""+key+"\"") {
																																																																																																																																																					t := strings.Split(s, "\"")
																																																																																																																																																								for i, t2 := range t {
																																																																																																																																																													if strings.Contains(t2, ":") {
																																																																																																																																																																			if t[i-1] == key {
																																																																																																																																																																										return t[i+1], nil
																																																																																																																																																																															}
																																																																																																																																																																																			}
																																																																																																																																																																																						}
																																																																																																																																																																																								}
																																																																																																																																																																																									}
																																																																																																																																																																																										line := strings.Split(s, ",")
																																																																																																																																																																																											for _, v := range line {
																																																																																																																																																																																														if strings.Contains(v, "\""+key+"\":") {
																																																																																																																																																																																																		t := strings.Split(v, "\"")
																																																																																																																																																																																																					for i, t2 := range t {
																																																																																																																																																																																																										if strings.Contains(t2, ":") {
																																																																																																																																																																																																																if t[i-1] == key {
																																																																																																																																																																																																																							return t[i+1], nil
																																																																																																																																																																																																																												}
																																																																																																																																																																																																																																}
																																																																																																																																																																																																																																			}
																																																																																																																																																																																																																																					}
																																																																																																																																																																																																																																						}
																																																																																																																																																																																																																																							return "", errors.Wrap(fmt.Errorf("can't find %s", key), "findKeyString")

																																																																																																																																																																																																																																						}


																																																																																																																																																																																																																																						func Loop1(wg *sync.WaitGroup)func(){
																																																																																																																																																																																																																																								return func() {
																																																																																																																																																																																																																																											defer fmt.Println("EXIT loop")
																																																																																																																																																																																																																																													defer wg.Done()
																																																																																																																																																																																																																																															fmt.Println("HAHA")
																																																																																																																																																																																																																																																	time.Sleep(10 * time.Second)
																																																																																																																																																																																																																																																			fmt.Println("_____")
																																																																																																																																																																																																																																																				}
																																																																																																																																																																																																																																																			}

																																																																																																																																																																																																																																																			func loop2(ctx context.Context)func(){
																																																																																																																																																																																																																																																					return func() {
																																																																																																																																																																																																																																																								//defer fmt.Println("exit loop")
																																																																																																																																																																																																																																																										fmt.Println("enter loop2")
																																																																																																																																																																																																																																																												time.Sleep(1 * time.Minute)
																																																																																																																																																																																																																																																														fmt.Println("**********")
																																																																																																																																																																																																																																																															}
																																																																																																																																																																																																																																																														}
