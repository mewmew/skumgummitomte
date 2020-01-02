package irgen

import (
	"fmt"
	gotypes "go/types"

	irtypes "github.com/llir/llvm/ir/types"
)

// ### [ Helper functions ] ####################################################

// RelStringer is the interface that wraps the Go SSA RelString method.
type RelStringer interface {
	// RelString returns the full name of the global, qualified by package name,
	// receiver type, etc.
	//
	// Examples:
	//
	//    "math.IsNaN"                  // a package-level function
	//    "(*bytes.Buffer).Bytes"       // a declared method or a wrapper
	//    "(*bytes.Buffer).Bytes$thunk" // thunk (func wrapping method; receiver is param 0)
	//    "(*bytes.Buffer).Bytes$bound" // bound (func wrapping method; receiver supplied by closure)
	//    "main.main$1"                 // an anonymous function in main
	//    "main.init#1"                 // a declared init function
	//    "main.init"                   // the synthesized package initializer
	RelString(from *gotypes.Package) string
}

// fullName returns the full name of the value, qualified by package name,
// receiver type, etc.
func (m *Module) fullName(v RelStringer) string {
	if m.goPkg.Pkg.Name() == "main" {
		// Fully qualified name if global is imported, otherwise name without
		// package path.
		from := m.goPkg.Pkg
		return v.RelString(from)
	}
	// Fully qualified name (with package path).
	return v.RelString(nil)
}

// precFromFloatKind return the precision of the given LLVM IR floating-point
// kind, where precision specifies the number of bits in the mantissa (including
// the implicit lead bit).
func precFromFloatKind(kind irtypes.FloatKind) uint {
	switch kind {
	// 16-bit floating-point type (IEEE 754 half precision).
	case irtypes.FloatKindHalf: // half
		return 11
	// 32-bit floating-point type (IEEE 754 single precision).
	case irtypes.FloatKindFloat: // float
		return 24
	// 64-bit floating-point type (IEEE 754 double precision).
	case irtypes.FloatKindDouble: // double
		return 53
	// 128-bit floating-point type (IEEE 754 quadruple precision).
	case irtypes.FloatKindFP128: // fp128
		return 113
	// 80-bit floating-point type (x86 extended precision).
	case irtypes.FloatKindX86_FP80: // x86_fp80
		return 64
	// 128-bit floating-point type (PowerPC double-double arithmetic).
	case irtypes.FloatKindPPC_FP128: // ppc_fp128
		return 106
	default:
		panic(fmt.Errorf("support for LLVM IR floating-point kind %v not yet implemented", kind))
	}
}
