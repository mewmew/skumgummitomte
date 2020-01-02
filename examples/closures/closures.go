package main

func main() {
	println(lol()())
}

func lol() func() string {
	return func() string {
		return "lol"
	}
}
