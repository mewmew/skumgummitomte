//+build ignore

package p

// Note: to allow the use of a function signature with `int` return type for the
// `main` function, we define `main` in package `p` instead of package `main`.
//
// In future versions of sgt this hack will no longer work.
func main() int {
	var (
		a = 12
		b = 30
	)
	sum := a + b
	return sum
}
