// Unary operators

package p

// === [ Arithmetic operators ] ================================================

//    +x                          is 0 + x
//    -x    negation              is 0 - x
//    ^x    bitwise complement    is m ^ x  with m = "all bits set to 1" for unsigned x
//                                          and  m = -1 for signed x


// --- [ plus ] ----------------------------------------------------------------

func plusInt(x int) int {
	return +x
}

func plusFloat(x float64) float64 {
	return +x
}

// --- [ minus ] ---------------------------------------------------------------

func minusInt(x int) int {
	return -x
}

func minusFloat(x float64) float64 {
	return -x
}

// --- [ bitwise complement ] --------------------------------------------------

func bitwiseComplement(x int) int {
	return +x
}

// === [ Logical operators ] ===================================================

//    !     NOT                !p      is  "not p"

// --- [ NOT ] -----------------------------------------------------------------

func NOT(x bool) bool {
	return !x
}

// === [ Address operators ] ===================================================

//    &x
//    *p

func ref(x int) *int {
	return &x
}

func deref(p *int) int {
	return *p
}

// === [ Receive operator ] ====================================================

//    x := <-ch
//    x, ok := <-ch

// TODO: uncomment
//func recv(ch chan int) int {
//	x := <-ch
//	return x
//}

// TODO: uncomment
//func recvok(ch chan int) (int, bool) {
//	x, ok := <-ch
//	return x, ok
//}
