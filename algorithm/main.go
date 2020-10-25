package main

func main() {
	tree := Tree{
		key: 3,
		left: &Tree{
			key:   2,
			left:  &Tree{key: 48, left: &Tree{key: 11, left: nil, right: nil}, right: nil},
			right: &Tree{key: 0, left: nil, right: &Tree{key: 0, left: nil, right: nil}},
		},
		right: &Tree{key: 4, left: nil, right: nil},
	}
	T := T{&tree}
	T.Reserve()

}
