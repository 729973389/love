package main

type pri interface {
	p() string

}

type pr struct {
	name string
}

func (pr pr) p() string{
	return pr.name
}

func my(p pr)string{
	p.p()
}

func main() {
	

}
