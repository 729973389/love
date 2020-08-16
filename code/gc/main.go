package main

import (
	"fmt"
)

const grey = 0;

const black = 1;

const white = 2;
type Mem struct {
	Color int
	Ptr *Mem
	Value string

}
func main () {
	m1 := &Mem{
		Color: grey,
		Ptr: &Mem{
			Color: white,
			Ptr: &Mem{
				Color: white,
				Ptr: nil,
				Value: "m1_3",
			},
			Value: "m1_2",
		},
		Value: "m1_1",
	}

	m2:=&Mem{
		Color: white,
		Ptr: nil,
		Value: "m2",
	}
	m3 :=&Mem{
		Color: grey,
		Ptr: m2,
		Value: "m3",
	}
	m4 := &Mem{
		Color: white,
		Ptr: &Mem{
			Color: white,
			Ptr: nil,
			Value: "m4_2",
		},
		Value: "m4_1",
	}
	var ma =make([]*Mem,0)
	ma=append(ma,m1,m2,m3,m4)
	for _,v:=range ma{
		Color(v)
	}
	for _,v := range ma {
		Gc(v)
		
	}


}

func Color(mem *Mem){
	if mem.Color==grey{
		mem.Color=black
		if mem.Ptr!=nil{
			mem.Ptr.Color=grey
			Color(mem.Ptr)
		}
	}
}

func Gc(mem *Mem){
	if mem.Color==white{
		fmt.Println("clean: ",mem.Value)
	}
	if mem.Ptr!=nil{
		Gc(mem.Ptr)
	}
}
