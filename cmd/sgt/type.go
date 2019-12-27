package main

import (
	"fmt"
	gotypes "go/types"

	irtypes "github.com/llir/llvm/ir/types"
	"golang.org/x/tools/go/ssa"
)

// initPredeclaredTypes initializes LLVM IR types corresponding to the
// predeclared types in Go (e.g. "bool").
func (gen *generator) initPredeclaredTypes() {
	// boolean type.
	boolType := irtypes.NewInt(1)
	boolType.SetName("bool")
	gen.predeclaredTypes[boolType.Name()] = boolType
	gen.module.TypeDefs = append(gen.module.TypeDefs, boolType)
	// signed integer types.
	intType := irtypes.NewInt(64) // default to 64-bit integer types.
	intType.SetName("int")
	gen.predeclaredTypes[intType.Name()] = intType
	gen.module.TypeDefs = append(gen.module.TypeDefs, intType)
	int8Type := irtypes.NewInt(8)
	int8Type.SetName("int8")
	gen.predeclaredTypes[int8Type.Name()] = int8Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, int8Type)
	int16Type := irtypes.NewInt(16)
	int16Type.SetName("int16")
	gen.predeclaredTypes[int16Type.Name()] = int16Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, int16Type)
	int32Type := irtypes.NewInt(32)
	int32Type.SetName("int32")
	gen.predeclaredTypes[int32Type.Name()] = int32Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, int32Type)
	int64Type := irtypes.NewInt(64)
	int64Type.SetName("int64")
	gen.predeclaredTypes[int64Type.Name()] = int64Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, int64Type)
	// unsigned integer types.
	uintType := irtypes.NewInt(64) // default to 64-bit integer types.
	uintType.SetName("uint")
	gen.predeclaredTypes[uintType.Name()] = uintType
	gen.module.TypeDefs = append(gen.module.TypeDefs, uintType)
	uint8Type := irtypes.NewInt(8)
	uint8Type.SetName("uint8")
	gen.predeclaredTypes[uint8Type.Name()] = uint8Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, uint8Type)
	uint16Type := irtypes.NewInt(16)
	uint16Type.SetName("uint16")
	gen.predeclaredTypes[uint16Type.Name()] = uint16Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, uint16Type)
	uint32Type := irtypes.NewInt(32)
	uint32Type.SetName("uint32")
	gen.predeclaredTypes[uint32Type.Name()] = uint32Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, uint32Type)
	uint64Type := irtypes.NewInt(64)
	uint64Type.SetName("uint64")
	gen.predeclaredTypes[uint64Type.Name()] = uint64Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, uint64Type)
	// unsigned integer pointer type.
	uintptrType := irtypes.NewInt(64) // default to 64-bit pointer types.
	uintptrType.SetName("uintptr")
	gen.predeclaredTypes[uintptrType.Name()] = uintptrType
	gen.module.TypeDefs = append(gen.module.TypeDefs, uintptrType)
	// floating-point types.
	float32Type := &irtypes.FloatType{Kind: irtypes.FloatKindFloat}
	float32Type.SetName("float32")
	gen.predeclaredTypes[float32Type.Name()] = float32Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, float32Type)
	float64Type := &irtypes.FloatType{Kind: irtypes.FloatKindDouble}
	float64Type.SetName("float64")
	gen.predeclaredTypes[float64Type.Name()] = float64Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, float64Type)
	// complex types.
	// TODO: add support for LLVM IR structure types with field names.
	//complex64Type = NewStruct(
	//   Field{Name: "real", Type: float32Type},
	//   Field{Name: "imag", Type: float32Type},
	//)
	complex64Type := irtypes.NewStruct(float32Type, float32Type)
	complex64Type.SetName("complex64")
	gen.predeclaredTypes[complex64Type.Name()] = complex64Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, complex64Type)
	// TODO: add support for LLVM IR structure types with field names.
	//complex128Type = NewStruct(
	//   Field{Name: "real", Type: float64Type},
	//   Field{Name: "imag", Type: float64Type},
	//)
	complex128Type := irtypes.NewStruct(float64Type, float64Type)
	complex128Type.SetName("complex128")
	gen.predeclaredTypes[complex128Type.Name()] = complex128Type
	gen.module.TypeDefs = append(gen.module.TypeDefs, complex128Type)
	// string type.
	// TODO: add support for LLVM IR structure types with field names.
	//stringType = NewStruct(
	//   Field{Name: "data", Type: irtypes.NewPointer(irtypes.NewInt(8))},
	//   Field{Name: "len", Type: intType},
	//)
	stringType := irtypes.NewStruct(irtypes.NewPointer(irtypes.NewInt(8)), intType)
	stringType.SetName("string")
	gen.predeclaredTypes[stringType.Name()] = stringType
	gen.module.TypeDefs = append(gen.module.TypeDefs, stringType)
	// unsafe pointer types.
	unsafePointerType := irtypes.NewPointer(irtypes.NewInt(8)) // void*
	unsafePointerType.SetName("unsafe.Pointer")
	gen.predeclaredTypes[unsafePointerType.Name()] = unsafePointerType
	gen.module.TypeDefs = append(gen.module.TypeDefs, unsafePointerType)
}

// llTypeFromName returns the LLVM IR type for the corresponding Go type name.
func (gen *generator) llTypeFromName(typeName string) irtypes.Type {
	// TODO: handle shadowing of predeclared types based on scope.
	if typ, ok := gen.predeclaredTypes[typeName]; ok {
		return typ
	}
	if typ, ok := gen.types[typeName]; ok {
		return typ
	}
	panic(fmt.Errorf("unable to locate LLVM IR type for the corresponding Go type name %q", typeName))
}

// compileType compiles the given Go SSA type into LLVM IR.
func (gen *generator) compileType(typeName string, goType *ssa.Type) error {
	// TODO: remove debug output.
	fmt.Println("compileType")
	fmt.Println(goType)
	fmt.Println()
	return nil
}

// llTypeFromGoType returns the LLVM IR type corresponding to the given Go type.
func (gen *generator) llTypeFromGoType(goType gotypes.Type) irtypes.Type {
	switch goType := goType.(type) {
	case *gotypes.Basic:
		return gen.llTypeFromGoBasicType(goType)
	case *gotypes.Pointer:
		return gen.llTypeFromGoPointerType(goType)
	default:
		panic(fmt.Errorf("support for Go type %T not yet implemented", goType))
	}
}

// llTypeFromGoBasicType returns the LLVM IR type corresponding to the given
// Go basic type.
func (gen *generator) llTypeFromGoBasicType(goType *gotypes.Basic) irtypes.Type {
	var typeName string
	switch kind := goType.Kind(); kind {
	// boolean type.
	case gotypes.Bool:
		typeName = "bool"
	// signed integer types.
	case gotypes.Int:
		typeName = "int"
	case gotypes.Int8:
		typeName = "int8"
	case gotypes.Int16:
		typeName = "int16"
	case gotypes.Int32:
		typeName = "int32"
	case gotypes.Int64:
		typeName = "int64"
	// unsigned integer types.
	case gotypes.Uint:
		typeName = "uint"
	case gotypes.Uint8:
		typeName = "uint8"
	case gotypes.Uint16:
		typeName = "uint16"
	case gotypes.Uint32:
		typeName = "uint32"
	case gotypes.Uint64:
		typeName = "uint64"
	// unsigned integer pointer type.
	case gotypes.Uintptr:
		typeName = "uintptr"
	// floating-point types.
	case gotypes.Float32:
		typeName = "float32"
	case gotypes.Float64:
		typeName = "float64"
	// complex types.
	case gotypes.Complex64:
		typeName = "complex64"
	case gotypes.Complex128:
		typeName = "complex128"
	// string type.
	case gotypes.String:
		typeName = "string"
	// unsafe pointer types.
	case gotypes.UnsafePointer:
		typeName = "unsafe.Pointer"
	// TODO: figure out how to handle untyped values if/when needed.
	// case gotypes.UntypedBool:
	// case gotypes.UntypedInt:
	// case gotypes.UntypedRune:
	// case gotypes.UntypedFloat:
	// case gotypes.UntypedComplex:
	// case gotypes.UntypedString:
	// case gotypes.UntypedNil:
	default:
		panic(fmt.Errorf("support for Go basic type kind %v not yet implemented", kind))
	}
	return gen.llTypeFromName(typeName)
}

// llTypeFromGoPointerType returns the LLVM IR type corresponding to the given
// Go pointer type.
func (gen *generator) llTypeFromGoPointerType(goType *gotypes.Pointer) *irtypes.PointerType {
	elemType := gen.llTypeFromGoType(goType.Elem())
	return irtypes.NewPointer(elemType)
}
