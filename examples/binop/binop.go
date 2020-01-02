// Binary operators

package p

// === [ Arithmetic operators ] ================================================
//
//    +    sum                    integers, floats, complex values, strings
//    -    difference             integers, floats, complex values
//    *    product                integers, floats, complex values
//    /    quotient               integers, floats, complex values
//    %    remainder              integers

// --- [ sum ] -----------------------------------------------------------------

func sumInt(x, y int) int {
	return x + y
}

func sumFloat(x, y float64) float64 {
	return x + y
}

func sumComplex(x, y complex128) complex128 {
	return x + y
}

// TODO: uncomment.
//func sumString(x, y string) string {
//	return x + y // string concatenation
//}

// --- [ difference ] ----------------------------------------------------------

func differenceInt(x, y int) int {
	return x - y
}

func differenceFloat(x, y float64) float64 {
	return x - y
}

func differenceComplex(x, y complex128) complex128 {
	return x - y
}

// --- [ product ] -------------------------------------------------------------

func productInt(x, y int) int {
	return x * y
}

func productFloat(x, y float64) float64 {
	return x * y
}

func productComplex(x, y complex128) complex128 {
	return x * y
}

// --- [ quotient ] ------------------------------------------------------------

func quotientInt(x, y int) int {
	return x / y
}

func quotientFloat(x, y float64) float64 {
	return x / y
}

func quotientComplex(x, y complex128) complex128 {
	return x / y
}

// --- [ remainder ] -----------------------------------------------------------

func remainder(x, y int) int {
	return x % y
}

//    &    bitwise AND            integers
//    |    bitwise OR             integers
//    ^    bitwise XOR            integers
//    &^   bit clear (AND NOT)    integers

// --- [ bitwise AND ] ---------------------------------------------------------

func bitwiseAND(x, y int) int {
	return x & y
}

// --- [ bitwise OR ] ----------------------------------------------------------

func bitwiseOR(x, y int) int {
	return x | y
}

// --- [ bitwise XOR ] ---------------------------------------------------------

func bitwiseXOR(x, y int) int {
	return x ^ y
}

// --- [ bit clear AND NOT ] ---------------------------------------------------

func bitClearANDNOT(x, y int) int {
	return x &^ y
}

//    <<   left shift             integer << unsigned integer
//    >>   right shift            integer >> unsigned integer

// --- [ left shift ] ----------------------------------------------------------

func leftShiftInt(x int, y uint) int {
	return x << y
}

func leftShiftUint(x, y uint) uint {
	return x << y
}

// --- [ right shift ] ---------------------------------------------------------

func rightShiftInt(x int, y uint) int {
	return x >> y
}

func rightShiftUint(x, y uint) uint {
	return x >> y
}

// === [ Comparison operators ] ================================================

// The equality operators == and != apply to operands that are comparable.
//
//    ==    equal
//    !=    not equal
//
// * Boolean values are comparable.
// * Integer values are comparable and ordered.
// * Floating-point values are comparable and ordered.
// * Complex values are comparable.
// * String values are comparable and ordered.
// * Pointer values are comparable.
// * Channel values are comparable.
// * Interface values are comparable.
// * Struct values are comparable if all their fields are comparable.
// * Array values are comparable if values of the array element type are
//   comparable.
// * Slice, map, and function values are not comparable (except against `nil`).

// --- [ equal ] ---------------------------------------------------------------

func equalBool(x, y bool) bool {
	return x == y
}

func equalInt(x, y int) bool {
	return x == y
}

func equalFloat(x, y float64) bool {
	return x == y
}

func equalComplex(x, y complex128) bool {
	return x == y
}

func equalString(x, y string) bool {
	return x == y
}

func equalPointer(x, y *int) bool {
	return x == y
}

// TODO: uncomment
//func equalChannel(x, y chan int) bool {
//	return x == y
//}

// TODO: uncomment
//func equalInterface(x, y interface{ Foo() int }) bool {
//	return x == y
//}

func equalStruct(x, y struct{ A int; B float64 }) bool {
	return x == y
}

func equalArray(x, y [10]int) bool {
	return x == y
}

func equalSliceNil(x []int) bool {
	return x == nil
}

// TODO: uncomment
//func equalMapNil(x map[int]float64) bool {
//	return x == nil
//}

func equalFuncNil(x func() int) bool {
	return x == nil
}

// --- [ not equal ] -----------------------------------------------------------

func notEqualBool(x, y bool) bool {
	return x != y
}

func notEqualInt(x, y int) bool {
	return x != y
}

func notEqualFloat(x, y float64) bool {
	return x != y
}

func notEqualComplex(x, y complex128) bool {
	return x != y
}

func notEqualString(x, y string) bool {
	return x != y
}

func notEqualPointer(x, y *int) bool {
	return x != y
}

// TODO: uncomment
//func notEqualChannel(x, y chan int) bool {
//	return x != y
//}

// TODO: uncomment
//func notEqualInterface(x, y interface{ Foo() int }) bool {
//	return x != y
//}

func notEqualStruct(x, y struct{ A int; B float64 }) bool {
	return x != y
}

func notEqualArray(x, y [10]int) bool {
	return x != y
}

func notEqualSliceNil(x []int) bool {
	return x != nil
}

// TODO: uncomment
//func notEqualMapNil(x map[int]float64) bool {
//	return x != nil
//}

func notEqualFuncNil(x func() int) bool {
	return x != nil
}

// The ordering operators <, <=, >, and >= apply to operands that are ordered.
//
//    <     less
//    <=    less or equal
//    >     greater
//    >=    greater or equal
//
// * Integer values are comparable and ordered.
// * Floating-point values are comparable and ordered.
// * String values are comparable and ordered.

// --- [ less ] ----------------------------------------------------------------

func lessInt(x, y int) bool {
	return x < y
}

func lessFloat(x, y float64) bool {
	return x < y
}

func lessString(x, y string) bool {
	return x < y
}

// --- [ less or equal ] -------------------------------------------------------

func lessOrEqualInt(x, y int) bool {
	return x <= y
}

func lessOrEqualFloat(x, y float64) bool {
	return x <= y
}

func lessOrEqualString(x, y string) bool {
	return x <= y
}

// --- [ greater ] -------------------------------------------------------------

func greaterInt(x, y int) bool {
	return x > y
}

func greaterFloat(x, y float64) bool {
	return x > y
}

func greaterString(x, y string) bool {
	return x > y
}

// --- [ greater or equal ] ----------------------------------------------------

func greaterOrEqualInt(x, y int) bool {
	return x >= y
}

func greaterOrEqualFloat(x, y float64) bool {
	return x >= y
}

func greaterOrEqualString(x, y string) bool {
	return x >= y
}

// === [ Logical operators ] ===================================================

//    &&    conditional AND    p && q  is  "if p then q else false"
//    ||    conditional OR     p || q  is  "if p then true else q"

// --- [ conditional AND ] -----------------------------------------------------

func conditionalAND(x, y bool) bool {
	return x && y
}

// --- [ conditional OR ] -----------------------------------------------------

func conditionalOR(x, y bool) bool {
	return x || y
}
