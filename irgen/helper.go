package irgen

import (
	"fmt"
	gotypes "go/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/metadata"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
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

// fullName returns the full name of the value, qualified by package name if not
// in main package.
func (m *Module) fullName(v RelStringer) string {
	if m.skipPkgPrefix() {
		// Fully qualified name if global is imported, otherwise name without
		// package path.
		from := m.goPkg.Pkg
		return v.RelString(from)
	}
	// Fully qualified name (with package path).
	return v.RelString(nil)
}

// fullTypeName returns the full name of the type, qualified by package name if
// not in main package.
func (m *Module) fullTypeName(t gotypes.Type) string {
	if m.skipPkgPrefix() {
		// Fully qualified name if type is imported, otherwise name without
		// package path.
		from := m.goPkg.Pkg
		return gotypes.TypeString(t, gotypes.RelativeTo(from))
	}
	// Fully qualified name (with package path).
	return gotypes.TypeString(t, nil)
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

// TODO: remove addMetadata when metadata attachment of llir/llvm instructions
// and terminators has been refined. Currently we have access to metadata
// attachments through the MDAttachments method, but cannot modify it.

// addMetadata adds the key-value pair as a metadata attachment to the given
// instruction or terminator.
func addMetadata(v ir.LLStringer, key, val string) {
	md := &metadata.Attachment{
		Name: key,
		Node: &metadata.Tuple{
			MetadataID: -1, // metadata literal.
			Fields: []metadata.Field{
				&metadata.String{Value: val},
			},
		},
	}
	switch v := v.(type) {
	// Unary instructions
	case *ir.InstFNeg:
		v.Metadata = append(v.Metadata, md)
	// Binary instructions
	case *ir.InstAdd:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFAdd:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstSub:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFSub:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstMul:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFMul:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstUDiv:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstSDiv:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFDiv:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstURem:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstSRem:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFRem:
		v.Metadata = append(v.Metadata, md)
	// Bitwise instructions
	case *ir.InstShl:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstLShr:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstAShr:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstAnd:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstOr:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstXor:
		v.Metadata = append(v.Metadata, md)
	// Vector instructions
	case *ir.InstExtractElement:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstInsertElement:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstShuffleVector:
		v.Metadata = append(v.Metadata, md)
	// Aggregate instructions
	case *ir.InstExtractValue:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstInsertValue:
		v.Metadata = append(v.Metadata, md)
	// Memory instructions
	case *ir.InstAlloca:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstLoad:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstStore:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFence:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstCmpXchg:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstAtomicRMW:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstGetElementPtr:
		v.Metadata = append(v.Metadata, md)
	// Conversion instructions
	case *ir.InstTrunc:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstZExt:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstSExt:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFPTrunc:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFPExt:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFPToUI:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFPToSI:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstUIToFP:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstSIToFP:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstPtrToInt:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstIntToPtr:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstBitCast:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstAddrSpaceCast:
		v.Metadata = append(v.Metadata, md)
	// Other instructions
	case *ir.InstICmp:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstFCmp:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstPhi:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstSelect:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstCall:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstVAArg:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstLandingPad:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstCatchPad:
		v.Metadata = append(v.Metadata, md)
	case *ir.InstCleanupPad:
		v.Metadata = append(v.Metadata, md)
	// Terminators
	case *ir.TermRet:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermBr:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermCondBr:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermSwitch:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermIndirectBr:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermInvoke:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermCallBr:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermResume:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermCatchSwitch:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermCatchRet:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermCleanupRet:
		v.Metadata = append(v.Metadata, md)
	case *ir.TermUnreachable:
		v.Metadata = append(v.Metadata, md)
	default:
		panic(fmt.Errorf("support for instruction or terminator %T not yet implemented", v))
	}
}

// isSigned reports whether the given integer type is signed.
func (m *Module) isSigned(typ *irtypes.IntType) bool {
	switch typ.Name() {
	case "uint8", "uint16", "uint32", "uint64", "uint", "uintptr":
		return false
	case "int8", "int16", "int32", "int64", "int":
		return true
	case "byte": // alias for uint8
		return false
	case "rune": // alias for int32
		return true
	default:
		panic(fmt.Errorf("support for integer type name %q not yet implemented", typ.Name()))
	}
}

// skipPkgPrefix reports whether to skip the package prefix in qualified names.
func (m *Module) skipPkgPrefix() bool {
	switch {
	case m.goPkg.Pkg.Name() == "main":
		// Skip package prefix for main packages.
		return true
	case m.goPkg.Pkg.Path() == "command-line-arguments":
		// Skip package prefix if a source file name was specified on command line
		// rather than a package path.
		return true
	default:
		return false
	}
}

// --- [ convert ] -------------------------------------------------------------

// convert converts the given the given value to the specified type, emitting to
// fn.
func (fn *Func) convert(from irvalue.Value, to irtypes.Type) irValueInstruction {
	var inst irValueInstruction
	switch fromType := from.Type().(type) {
	case *irtypes.IntType:
		switch toType := to.(type) {
		// int -> int
		case *irtypes.IntType:
			switch {
			case fromType.BitSize == toType.BitSize:
				inst = fn.cur.NewBitCast(from, to)
			case fromType.BitSize < toType.BitSize:
				if fn.m.isSigned(fromType) {
					inst = fn.cur.NewSExt(from, to)
				} else {
					inst = fn.cur.NewZExt(from, to)
				}
			case fromType.BitSize > toType.BitSize:
				inst = fn.cur.NewTrunc(from, to)
			default:
				panic(fmt.Errorf("support for converting from type %T (%v) to type %T (%v) not yet implemented", fromType, fromType, to, to))
			}
		// int -> float
		case *irtypes.FloatType:
			if fn.m.isSigned(fromType) {
				inst = fn.cur.NewSIToFP(from, to)
			} else {
				inst = fn.cur.NewUIToFP(from, to)
			}
		// TODO: add support for more to types.
		default:
			panic(fmt.Errorf("support for converting from type %T (%v) to type %T (%v) not yet implemented", fromType, fromType, to, to))
		}
	case *irtypes.FloatType:
		switch toType := to.(type) {
		// float -> int
		case *irtypes.IntType:
			if fn.m.isSigned(toType) {
				inst = fn.cur.NewFPToSI(from, to)
			} else {
				inst = fn.cur.NewFPToUI(from, to)
			}
		// float -> float
		case *irtypes.FloatType:
			fromPrec := precFromFloatKind(fromType.Kind)
			toPrec := precFromFloatKind(toType.Kind)
			switch {
			case fromPrec == toPrec:
				// Currently the number of precision bits is unique to each
				// floating-point kind, thus this simply converts from e.g. double
				// to double. Otherwise, bitcast would not be valid.
				// TODO: consider using FPTrunc of FPExt instead, as they are
				// float-aware (which bitcast is not).
				inst = fn.cur.NewBitCast(from, to)
			case fromPrec < toPrec:
				inst = fn.cur.NewFPExt(from, to)
			case fromPrec > toPrec:
				inst = fn.cur.NewFPTrunc(from, to)
			}
		// TODO: add support for more to types.
		default:
			panic(fmt.Errorf("support for converting from type %T (%v) to type %T (%v) not yet implemented", fromType, fromType, to, to))
		}
	case *irtypes.PointerType:
		switch to.(type) {
		// pointer -> pointer
		case *irtypes.PointerType:
			inst = fn.cur.NewBitCast(from, to)
		// TODO: add support for more to types.
		default:
			panic(fmt.Errorf("support for converting from type %T (%v) to type %T (%v) not yet implemented", fromType, fromType, to, to))
		}
	// TODO: add support for more from types.
	default:
		panic(fmt.Errorf("support for converting from type %T (%v) to type %T (%v) not yet implemented", fromType, fromType, to, to))
	}
	return inst
}
