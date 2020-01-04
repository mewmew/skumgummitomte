package main

func main() {
	s := []string{
		"foo",
		"bar",
		"baz",
		"qux",
	}
	t := f(s)
	println(t[0])
	println(t[1])
	println(t[2])
	println(t[3])
}

func f(s []string) []string {
	s[0], s[1], s[2], s[3] = s[3], s[2], s[1], s[0]
	return s
}

func setlow() {
	var a []int
	_ = a[1:]
}

func sethigh() {
	var a []int
	_ = a[:2]
}

func setmax() {
	var a []int
	_ = a[:2:3]
}

func slicedef() {
	var a []int = []int{1, 2, 3}
	_ = a
}
