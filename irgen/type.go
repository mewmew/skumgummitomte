package irgen

import (
	"fmt"
	gotypes "go/types"

	irtypes "github.com/llir/llvm/ir/types"
	"golang.org/x/tools/go/ssa"
)

// --- [ init ] ----------------------------------------------------------------

// initPredeclaredTypes initializes LLVM IR types corresponding to the
// predeclared types in Go (e.g. "bool").
func (m *Module) initPredeclaredTypes() {
	// boolean type.
	boolType := irtypes.NewInt(1)
	boolType.SetName("bool")
	m.predeclaredTypes[boolType.Name()] = boolType
	m.Module.TypeDefs = append(m.Module.TypeDefs, boolType)
	// signed integer types.
	intType := irtypes.NewInt(64) // default to 64-bit integer types.
	intType.SetName("int")
	m.predeclaredTypes[intType.Name()] = intType
	m.Module.TypeDefs = append(m.Module.TypeDefs, intType)
	int8Type := irtypes.NewInt(8)
	int8Type.SetName("int8")
	m.predeclaredTypes[int8Type.Name()] = int8Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int8Type)
	int16Type := irtypes.NewInt(16)
	int16Type.SetName("int16")
	m.predeclaredTypes[int16Type.Name()] = int16Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int16Type)
	int32Type := irtypes.NewInt(32)
	int32Type.SetName("int32")
	m.predeclaredTypes[int32Type.Name()] = int32Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int32Type)
	int64Type := irtypes.NewInt(64)
	int64Type.SetName("int64")
	m.predeclaredTypes[int64Type.Name()] = int64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int64Type)
	// unsigned integer types.
	uintType := irtypes.NewInt(64) // default to 64-bit integer types.
	uintType.SetName("uint")
	m.predeclaredTypes[uintType.Name()] = uintType
	m.Module.TypeDefs = append(m.Module.TypeDefs, uintType)
	uint8Type := irtypes.NewInt(8)
	uint8Type.SetName("uint8")
	m.predeclaredTypes[uint8Type.Name()] = uint8Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint8Type)
	uint16Type := irtypes.NewInt(16)
	uint16Type.SetName("uint16")
	m.predeclaredTypes[uint16Type.Name()] = uint16Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint16Type)
	uint32Type := irtypes.NewInt(32)
	uint32Type.SetName("uint32")
	m.predeclaredTypes[uint32Type.Name()] = uint32Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint32Type)
	uint64Type := irtypes.NewInt(64)
	uint64Type.SetName("uint64")
	m.predeclaredTypes[uint64Type.Name()] = uint64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint64Type)
	// unsigned integer pointer type.
	uintptrType := irtypes.NewInt(64) // default to 64-bit pointer types.
	uintptrType.SetName("uintptr")
	m.predeclaredTypes[uintptrType.Name()] = uintptrType
	m.Module.TypeDefs = append(m.Module.TypeDefs, uintptrType)
	// floating-point types.
	float32Type := &irtypes.FloatType{Kind: irtypes.FloatKindFloat}
	float32Type.SetName("float32")
	m.predeclaredTypes[float32Type.Name()] = float32Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, float32Type)
	float64Type := &irtypes.FloatType{Kind: irtypes.FloatKindDouble}
	float64Type.SetName("float64")
	m.predeclaredTypes[float64Type.Name()] = float64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, float64Type)
	// complex types.
	// TODO: add support for LLVM IR structure types with field names.
	//complex64Type = NewStruct(
	//   Field{Name: "real", Type: float32Type},
	//   Field{Name: "imag", Type: float32Type},
	//)
	complex64Type := irtypes.NewStruct(float32Type, float32Type)
	complex64Type.SetName("complex64")
	m.predeclaredTypes[complex64Type.Name()] = complex64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, complex64Type)
	// TODO: add support for LLVM IR structure types with field names.
	//complex128Type = NewStruct(
	//   Field{Name: "real", Type: float64Type},
	//   Field{Name: "imag", Type: float64Type},
	//)
	complex128Type := irtypes.NewStruct(float64Type, float64Type)
	complex128Type.SetName("complex128")
	m.predeclaredTypes[complex128Type.Name()] = complex128Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, complex128Type)
	// string type.
	// TODO: add support for LLVM IR structure types with field names.
	//stringType = NewStruct(
	//   Field{Name: "data", Type: irtypes.NewPointer(irtypes.NewInt(8))},
	//   Field{Name: "len", Type: intType},
	//)
	stringType := irtypes.NewStruct(irtypes.NewPointer(irtypes.NewInt(8)), intType)
	stringType.SetName("string")
	m.predeclaredTypes[stringType.Name()] = stringType
	m.Module.TypeDefs = append(m.Module.TypeDefs, stringType)
	// unsafe pointer types.
	unsafePointerType := irtypes.NewPointer(irtypes.NewInt(8)) // void*
	unsafePointerType.SetName("unsafe.Pointer")
	m.predeclaredTypes[unsafePointerType.Name()] = unsafePointerType
	m.Module.TypeDefs = append(m.Module.TypeDefs, unsafePointerType)
}

// --- [ get ] -----------------------------------------------------------------

// irTypeFromName returns the LLVM IR type for the corresponding Go type name.
func (m *Module) irTypeFromName(typeName string) irtypes.Type {
	// TODO: handle shadowing of predeclared types based on scope.
	if typ, ok := m.predeclaredTypes[typeName]; ok {
		return typ
	}
	if typ, ok := m.types[typeName]; ok {
		return typ
	}
	panic(fmt.Errorf("unable to locate LLVM IR type of Go type with name %q", typeName))
}

// --- [ convert ] -------------------------------------------------------------

// TODO: check if precondition is needed. Type index probably needed for cyclic
// types.

// irTypeFromGo returns the LLVM IR type corresponding to the given Go type,
// emitting to m.
//
// Pre-condition: index named types.
func (m *Module) irTypeFromGo(goType gotypes.Type) irtypes.Type {
	switch goType := goType.(type) {
	case *gotypes.Array:
		panic("support for *gotypes.Array not yet implemented")
	case *gotypes.Basic:
		return m.irTypeFromGoBasicType(goType)
	case *gotypes.Chan:
		panic("support for *gotypes.Chan not yet implemented")
	case *gotypes.Interface:
		panic("support for *gotypes.Interface not yet implemented")
	case *gotypes.Map:
		panic("support for *gotypes.Map not yet implemented")
	case *gotypes.Named:
		return m.irTypeFromName(goType.Obj().Id())
	case *gotypes.Pointer:
		return m.irTypeFromGoPointerType(goType)
	case *gotypes.Signature:
		panic("support for *gotypes.Signature not yet implemented")
	case *gotypes.Slice:
		return m.irTypeFromGoSliceType(goType)
	case *gotypes.Struct:
		return m.irTypeFromGoStructType(goType)
	case *gotypes.Tuple:
		panic("support for *gotypes.Tuple not yet implemented")
	default:
		panic(fmt.Errorf("support for Go type %T not yet implemented", goType))
	}
}

// ~~~ [ basic type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoBasicType returns the LLVM IR type corresponding to the given Go
// basic type, emitting to m.
func (m *Module) irTypeFromGoBasicType(goType *gotypes.Basic) irtypes.Type {
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
	return m.irTypeFromName(typeName)
}

// ~~~ [ slice type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoSliceType returns the LLVM IR type corresponding to the given Go
// slice type, emitting to m.
func (m *Module) irTypeFromGoSliceType(goType *gotypes.Slice) *irtypes.StructType {
	elemType := m.irTypeFromGo(goType.Elem())
	data := irtypes.NewPointer(elemType)
	length := m.irTypeFromName("int")
	capacity := m.irTypeFromName("int")
	return irtypes.NewStruct(data, length, capacity)
}

// ~~~ [ struct type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoStructType returns the LLVM IR type corresponding to the given Go
// struct type, emitting to m.
func (m *Module) irTypeFromGoStructType(goType *gotypes.Struct) *irtypes.StructType {
	var fields []irtypes.Type
	for i := 0; i < goType.NumFields(); i++ {
		goField := goType.Field(i)
		if goField.Embedded() {
			// TODO: add support for embedded types.
			panic(fmt.Errorf("support for embedded types not yet implemented; struct type has embedded field %q", goField.Name()))
		}
		// TODO: add custom LLVM IR struct type which retains the names of struct
		// fields.
		name := goField.Name()
		_ = name
		field := m.irTypeFromGo(goField.Type())
		fields = append(fields, field)
	}
	return irtypes.NewStruct(fields...)
}

// ~~~ [ pointer type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoPointerType returns the LLVM IR type corresponding to the given
// Go pointer type, emitting to m.
func (m *Module) irTypeFromGoPointerType(goType *gotypes.Pointer) *irtypes.PointerType {
	elemType := m.irTypeFromGo(goType.Elem())
	return irtypes.NewPointer(elemType)
}

// --- [ index ] ---------------------------------------------------------------

// TODO: remove indexType?

// indexType indexes the given Go SSA type, creating a corresponding LLVM IR
// type, emitting to m.
func (m *Module) indexType(goType *ssa.Type) error {
	panic("not yet implemented")
}

// --- [ compile ] -------------------------------------------------------------

// emitType compiles the given Go SSA type definition to corresponding LLVM IR,
// emitting to m.
func (m *Module) emitType(goType *ssa.Type) error {
	dbg.Println("emitType")
	underlying := m.irTypeFromGo(goType.Type().Underlying())
	// Perform a deep copy of the underlying type. Otherwise, we may reset the
	// name of a previously named type.
	// TODO: only deep copy underlying type if it is a named type. Also, consider
	// doing a shallow copy instead of deep copy, as that should be enough to set
	// the type definition name (with the risk of resetting the name of a
	// previous type definition), but allows sharing the underlying types.
	typ := m.copyType(underlying)
	typ.SetName(goType.Name())
	dbg.Printf("   typ: %s = type %s", typ.String(), typ.LLString())
	m.types[typ.Name()] = typ
	m.TypeDefs = append(m.TypeDefs, typ)
	return nil
}

// --- [ copy ] ----------------------------------------------------------------

// copyType returns an identical copy of the given LLVM IR type.
func (m *Module) copyType(t irtypes.Type) irtypes.Type {
	switch t := t.(type) {
	case *irtypes.VoidType:
		return m.copyVoidType(t)
	case *irtypes.FuncType:
		return m.copyFuncType(t)
	case *irtypes.IntType:
		return m.copyIntType(t)
	case *irtypes.FloatType:
		return m.copyFloatType(t)
	case *irtypes.MMXType:
		return m.copyMMXType(t)
	case *irtypes.PointerType:
		return m.copyPointerType(t)
	case *irtypes.VectorType:
		return m.copyVectorType(t)
	case *irtypes.LabelType:
		return m.copyLabelType(t)
	case *irtypes.TokenType:
		return m.copyTokenType(t)
	case *irtypes.MetadataType:
		return m.copyMetadataType(t)
	case *irtypes.ArrayType:
		return m.copyArrayType(t)
	case *irtypes.StructType:
		return m.copyStructType(t)
	default:
		panic(fmt.Errorf("support for %T not yet implemented", t))
	}
}

// copyVoidType returns an identical copy of the given LLVM IR type.
func (m *Module) copyVoidType(t *irtypes.VoidType) *irtypes.VoidType {
	return &irtypes.VoidType{
		TypeName: t.TypeName,
	}
}

// copyFuncType returns an identical copy of the given LLVM IR type.
func (m *Module) copyFuncType(t *irtypes.FuncType) *irtypes.FuncType {
	var params []irtypes.Type
	for i := range t.Params {
		param := m.copyType(t.Params[i])
		params = append(params, param)
	}
	return &irtypes.FuncType{
		TypeName: t.TypeName,
		RetType:  m.copyType(t.RetType),
		Params:   params,
		Variadic: t.Variadic,
	}
}

// copyIntType returns an identical copy of the given LLVM IR type.
func (m *Module) copyIntType(t *irtypes.IntType) *irtypes.IntType {
	return &irtypes.IntType{
		TypeName: t.TypeName,
		BitSize:  t.BitSize,
	}
}

// copyFloatType returns an identical copy of the given LLVM IR type.
func (m *Module) copyFloatType(t *irtypes.FloatType) *irtypes.FloatType {
	return &irtypes.FloatType{
		TypeName: t.TypeName,
		Kind:     t.Kind,
	}
}

// copyMMXType returns an identical copy of the given LLVM IR type.
func (m *Module) copyMMXType(t *irtypes.MMXType) *irtypes.MMXType {
	return &irtypes.MMXType{
		TypeName: t.TypeName,
	}
}

// copyPointerType returns an identical copy of the given LLVM IR type.
func (m *Module) copyPointerType(t *irtypes.PointerType) *irtypes.PointerType {
	return &irtypes.PointerType{
		TypeName:  t.TypeName,
		ElemType:  m.copyType(t.ElemType),
		AddrSpace: t.AddrSpace,
	}
}

// copyVectorType returns an identical copy of the given LLVM IR type.
func (m *Module) copyVectorType(t *irtypes.VectorType) *irtypes.VectorType {
	return &irtypes.VectorType{
		TypeName: t.TypeName,
		Scalable: t.Scalable,
		Len:      t.Len,
		ElemType: m.copyType(t.ElemType),
	}
}

// copyLabelType returns an identical copy of the given LLVM IR type.
func (m *Module) copyLabelType(t *irtypes.LabelType) *irtypes.LabelType {
	return &irtypes.LabelType{
		TypeName: t.TypeName,
	}
}

// copyTokenType returns an identical copy of the given LLVM IR type.
func (m *Module) copyTokenType(t *irtypes.TokenType) *irtypes.TokenType {
	return &irtypes.TokenType{
		TypeName: t.TypeName,
	}
}

// copyMetadataType returns an identical copy of the given LLVM IR type.
func (m *Module) copyMetadataType(t *irtypes.MetadataType) *irtypes.MetadataType {
	return &irtypes.MetadataType{
		TypeName: t.TypeName,
	}
}

// copyArrayType returns an identical copy of the given LLVM IR type.
func (m *Module) copyArrayType(t *irtypes.ArrayType) *irtypes.ArrayType {
	return &irtypes.ArrayType{
		TypeName: t.TypeName,
		Len:      t.Len,
		ElemType: m.copyType(t.ElemType),
	}
}

// copyStructType returns an identical copy of the given LLVM IR type.
func (m *Module) copyStructType(t *irtypes.StructType) *irtypes.StructType {
	var fields []irtypes.Type
	for i := range t.Fields {
		field := m.copyType(t.Fields[i])
		fields = append(fields, field)
	}
	return &irtypes.StructType{
		TypeName: t.TypeName,
		Packed:   t.Packed,
		Fields:   fields,
		Opaque:   t.Opaque,
	}
}
