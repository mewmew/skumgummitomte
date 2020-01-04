package main

func main() {
	cmp("foo", "foo")
	cmp("abc", "def")
	cmp("abc", "abcd")
	cmp("abcd", "abc")
	cmp("abx", "abc")
}

func cmp(x, y string) {
	println("x:")
	println(x)
	println("y:")
	println(y)
	if x < y {
		println("x less than y")
	}
	if x > y {
		println("x greater than y")
	}
	if x <= y {
		println("x less than or equal to y")
	}
	if x >= y {
		println("x greater than or equal to y")
	}
	if x != y {
		println("x not equal to y")
	}
	if x == y {
		println("x equal to y")
	}
	println("")
}
