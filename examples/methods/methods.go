package main

type T int

func main() {
	var t T
	t.M1()
	t.M2()
}

func (t T) M1() {
	println("T.M1")
}

func (t *T) M2() {
	println("T.M2")
}
