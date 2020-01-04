package p

// Note: to allow the use of a function signature with `int` return type for the
// `main` function, we define `main` in package `p` instead of package `main`.
//
// In future versions of sgt this hack will no longer work.
func main() int {
	s := "The quick brown fox jumps over a lazy dog."
	return len(s)
}
