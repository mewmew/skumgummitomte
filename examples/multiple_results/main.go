package main

func main() {
	foo, bar := swap("foo", "bar")
	println(foo)
	println(bar)
}

func swap(a, b string) (string, string) {
	return b, a
}
