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
	m.types[boolType.Name()] = boolType
	m.Module.TypeDefs = append(m.Module.TypeDefs, boolType)
	// signed integer types.
	intType := irtypes.NewInt(64) // default to 64-bit integer types.
	intType.SetName("int")
	m.types[intType.Name()] = intType
	m.Module.TypeDefs = append(m.Module.TypeDefs, intType)
	int8Type := irtypes.NewInt(8)
	int8Type.SetName("int8")
	m.types[int8Type.Name()] = int8Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int8Type)
	int16Type := irtypes.NewInt(16)
	int16Type.SetName("int16")
	m.types[int16Type.Name()] = int16Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int16Type)
	int32Type := irtypes.NewInt(32)
	int32Type.SetName("int32")
	m.types[int32Type.Name()] = int32Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int32Type)
	int64Type := irtypes.NewInt(64)
	int64Type.SetName("int64")
	m.types[int64Type.Name()] = int64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, int64Type)
	// unsigned integer types.
	uintType := irtypes.NewInt(64) // default to 64-bit integer types.
	uintType.SetName("uint")
	m.types[uintType.Name()] = uintType
	m.Module.TypeDefs = append(m.Module.TypeDefs, uintType)
	uint8Type := irtypes.NewInt(8)
	uint8Type.SetName("uint8")
	m.types[uint8Type.Name()] = uint8Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint8Type)
	uint16Type := irtypes.NewInt(16)
	uint16Type.SetName("uint16")
	m.types[uint16Type.Name()] = uint16Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint16Type)
	uint32Type := irtypes.NewInt(32)
	uint32Type.SetName("uint32")
	m.types[uint32Type.Name()] = uint32Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint32Type)
	uint64Type := irtypes.NewInt(64)
	uint64Type.SetName("uint64")
	m.types[uint64Type.Name()] = uint64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, uint64Type)
	// unsigned integer pointer type.
	uintptrType := irtypes.NewInt(64) // default to 64-bit pointer types.
	uintptrType.SetName("uintptr")
	m.types[uintptrType.Name()] = uintptrType
	m.Module.TypeDefs = append(m.Module.TypeDefs, uintptrType)
	// floating-point types.
	float32Type := &irtypes.FloatType{Kind: irtypes.FloatKindFloat}
	float32Type.SetName("float32")
	m.types[float32Type.Name()] = float32Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, float32Type)
	float64Type := &irtypes.FloatType{Kind: irtypes.FloatKindDouble}
	float64Type.SetName("float64")
	m.types[float64Type.Name()] = float64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, float64Type)
	// complex types.
	// TODO: add support for LLVM IR structure types with field names.
	//complex64Type = NewStruct(
	//   Field{Name: "real", Type: float32Type},
	//   Field{Name: "imag", Type: float32Type},
	//)
	complex64Type := irtypes.NewStruct(float32Type, float32Type)
	complex64Type.SetName("complex64")
	m.types[complex64Type.Name()] = complex64Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, complex64Type)
	// TODO: add support for LLVM IR structure types with field names.
	//complex128Type = NewStruct(
	//   Field{Name: "real", Type: float64Type},
	//   Field{Name: "imag", Type: float64Type},
	//)
	complex128Type := irtypes.NewStruct(float64Type, float64Type)
	complex128Type.SetName("complex128")
	m.types[complex128Type.Name()] = complex128Type
	m.Module.TypeDefs = append(m.Module.TypeDefs, complex128Type)
	// string type.
	// TODO: add support for LLVM IR structure types with field names.
	//stringType = NewStruct(
	//   Field{Name: "data", Type: irtypes.NewPointer(uint8Type)}, // TODO: use byte alias instead of uint8.
	//   Field{Name: "len", Type: intType},
	//)
	stringType := irtypes.NewStruct(irtypes.NewPointer(uint8Type), intType) // TODO: use byte alias instead of uint8.
	stringType.SetName("string")
	m.types[stringType.Name()] = stringType
	m.Module.TypeDefs = append(m.Module.TypeDefs, stringType)
	// unsafe pointer type.
	unsafePointerType := irtypes.NewPointer(irtypes.NewInt(8)) // void*
	unsafePointerType.SetName("unsafe.Pointer")
	m.types[unsafePointerType.Name()] = unsafePointerType
	m.Module.TypeDefs = append(m.Module.TypeDefs, unsafePointerType)
	// error interface type.
	// TODO: figure out how to define interface types.
	// TODO: add support for LLVM IR structure types with field names.
	//errorType = NewStruct(
	//   Field{Name: "type", Type: stringType},
	//   Field{Name: "value", Type: irtypes.I8Ptr}, // generic pointer type
	//)
	errorType := irtypes.NewStruct(
		stringType,
		irtypes.I8Ptr, // generic pointer type
	)
	errorType.SetName("error")
	m.types[errorType.Name()] = errorType
	m.Module.TypeDefs = append(m.Module.TypeDefs, errorType)
}

// --- [ get ] -----------------------------------------------------------------

// irTypeFromName returns the LLVM IR type for the corresponding Go type name.
func (m *Module) irTypeFromName(typeName string) irtypes.Type {
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
		return m.irTypeFromGoArrayType(goType)
	case *gotypes.Basic:
		return m.irTypeFromGoBasicType(goType)
	case *gotypes.Chan:
		panic("support for *gotypes.Chan not yet implemented")
	case *gotypes.Interface:
		return m.irTypeFromGoInterfaceType(goType)
	case *gotypes.Map:
		panic("support for *gotypes.Map not yet implemented")
	case *gotypes.Named:
		typeName := m.fullTypeName(goType)
		return m.irTypeFromName(typeName)
	case *gotypes.Pointer:
		return m.irTypeFromGoPointerType(goType)
	case *gotypes.Signature:
		return m.irTypeFromGoSignatureType(goType)
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

// ~~~ [ array type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoArrayType returns the LLVM IR type corresponding to the given Go
// array type, emitting to m.
func (m *Module) irTypeFromGoArrayType(goType *gotypes.Array) *irtypes.ArrayType {
	length := uint64(goType.Len())
	elemType := m.irTypeFromGo(goType.Elem())
	return irtypes.NewArray(length, elemType)
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
	case gotypes.UntypedBool:
		typeName = "bool"
	case gotypes.UntypedInt:
		typeName = "int"
	// case gotypes.UntypedRune:
	case gotypes.UntypedFloat:
		typeName = "float64"
	// case gotypes.UntypedComplex:
	case gotypes.UntypedString:
		typeName = "string"
	// case gotypes.UntypedNil:
	default:
		panic(fmt.Errorf("support for Go basic type kind %v not yet implemented", kind))
	}
	return m.irTypeFromName(typeName)
}

// ~~~ [ interface type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoInterfaceType returns the LLVM IR type corresponding to the given
// Go interface type, emitting to m.
func (m *Module) irTypeFromGoInterfaceType(goType *gotypes.Interface) irtypes.Type {
	panic("not yet implemented")
}

// ~~~ [ pointer type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoPointerType returns the LLVM IR type corresponding to the given
// Go pointer type, emitting to m.
func (m *Module) irTypeFromGoPointerType(goType *gotypes.Pointer) *irtypes.PointerType {
	elemType := m.irTypeFromGo(goType.Elem())
	return irtypes.NewPointer(elemType)
}

// ~~~ [ function signature type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoSignatureType returns the LLVM IR type corresponding to the given
// Go function signature type, emitting to m.
func (m *Module) irTypeFromGoSignatureType(goType *gotypes.Signature) *irtypes.PointerType { //*irtypes.FuncType {
	if goType.Recv() != nil {
		// TODO: add support for methods; add receiver as first function
		// parameter.
		panic("support for methods in Go function signature type not yet implemented")
	}
	// Convert Go function parameters to equivalent LLVM IR function parameter
	// types.
	var paramTypes []irtypes.Type
	goParams := goType.Params()
	for i := 0; i < goParams.Len(); i++ {
		goParam := goParams.At(i)
		paramName := goParam.Name()
		_ = paramName // TODO: add parameter name as metadata?
		paramType := m.irTypeFromGo(goParam.Type())
		paramTypes = append(paramTypes, paramType)
	}
	// Convert Go function return types to equivalent LLVM IR function return
	// types.
	var resultTypes []irtypes.Type
	goResults := goType.Results()
	for i := 0; i < goResults.Len(); i++ {
		goResult := goResults.At(i)
		resultName := goResult.Name()
		// TODO: add resultName as field name of (custom) result structure type.
		_ = resultName
		resultType := m.irTypeFromGo(goResult.Type())
		resultTypes = append(resultTypes, resultType)
	}
	// Convert multiple return types a single return type by creating a structure
	// type with one field per return type.
	var retType irtypes.Type
	switch len(resultTypes) {
	// void return.
	case 0:
		retType = irtypes.Void
	// single return type.
	case 1:
		retType = resultTypes[0]
	// multiple return types.
	default:
		retType = irtypes.NewStruct(resultTypes...)
	}
	// Generate LLVM IR function signature type.
	sig := irtypes.NewFunc(retType, paramTypes...)
	sig.Variadic = goType.Variadic()
	// TODO: consider how to convert Go function signature types, in LLVM IR,
	// function values when used as callees are of type pointer to function
	// signature type.
	return irtypes.NewPointer(sig)
}

// ~~~ [ slice type ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irTypeFromGoSliceType returns the LLVM IR type corresponding to the given Go
// slice type, emitting to m.
func (m *Module) irTypeFromGoSliceType(goType *gotypes.Slice) *irtypes.StructType {
	elemType := m.irTypeFromGo(goType.Elem())
	return m.newSliceType(elemType)
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
	typeName := m.fullName(goType)
	dbg.Println("   typeName:", typeName)
	if _, ok := m.types[typeName]; ok {
		// type definition already present.
		return nil
	}
	underlying := m.irTypeFromGo(goType.Type().Underlying())
	// Perform a deep copy of the underlying type. Otherwise, we may reset the
	// name of a previously named type.
	// TODO: only deep copy underlying type if it is a named type. Also, consider
	// doing a shallow copy instead of deep copy, as that should be enough to set
	// the type definition name (with the risk of resetting the name of a
	// previous type definition), but allows sharing the underlying types.
	typ := copyTypeShallow(underlying)
	typ.SetName(typeName)
	dbg.Printf("   typ: %s = type %s", typ.String(), typ.LLString())
	m.types[typ.Name()] = typ
	m.TypeDefs = append(m.TypeDefs, typ)
	return nil
}

// --- [ copy ] ----------------------------------------------------------------

// ~~~ [ shallow copy ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// copyTypeShallow returns an identical copy of the given LLVM IR type.
func copyTypeShallow(t irtypes.Type) irtypes.Type {
	switch t := t.(type) {
	case *irtypes.VoidType:
		u := &irtypes.VoidType{}
		*u = *t
		return u
	case *irtypes.FuncType:
		u := &irtypes.FuncType{}
		*u = *t
		return u
	case *irtypes.IntType:
		u := &irtypes.IntType{}
		*u = *t
		return u
	case *irtypes.FloatType:
		u := &irtypes.FloatType{}
		*u = *t
		return u
	case *irtypes.MMXType:
		u := &irtypes.MMXType{}
		*u = *t
		return u
	case *irtypes.PointerType:
		u := &irtypes.PointerType{}
		*u = *t
		return u
	case *irtypes.VectorType:
		u := &irtypes.VectorType{}
		*u = *t
		return u
	case *irtypes.LabelType:
		u := &irtypes.LabelType{}
		*u = *t
		return u
	case *irtypes.TokenType:
		u := &irtypes.TokenType{}
		*u = *t
		return u
	case *irtypes.MetadataType:
		u := &irtypes.MetadataType{}
		*u = *t
		return u
	case *irtypes.ArrayType:
		u := &irtypes.ArrayType{}
		*u = *t
		return u
	case *irtypes.StructType:
		u := &irtypes.StructType{}
		*u = *t
		return u
	default:
		panic(fmt.Errorf("support for %T not yet implemented", t))
	}
}

// ~~~ [ deep copy ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// copyType returns an identical copy of the given LLVM IR type.
func copyType(t irtypes.Type) irtypes.Type {
	switch t := t.(type) {
	case *irtypes.VoidType:
		return copyVoidType(t)
	case *irtypes.FuncType:
		return copyFuncType(t)
	case *irtypes.IntType:
		return copyIntType(t)
	case *irtypes.FloatType:
		return copyFloatType(t)
	case *irtypes.MMXType:
		return copyMMXType(t)
	case *irtypes.PointerType:
		return copyPointerType(t)
	case *irtypes.VectorType:
		return copyVectorType(t)
	case *irtypes.LabelType:
		return copyLabelType(t)
	case *irtypes.TokenType:
		return copyTokenType(t)
	case *irtypes.MetadataType:
		return copyMetadataType(t)
	case *irtypes.ArrayType:
		return copyArrayType(t)
	case *irtypes.StructType:
		return copyStructType(t)
	default:
		panic(fmt.Errorf("support for %T not yet implemented", t))
	}
}

// copyVoidType returns an identical copy of the given LLVM IR type.
func copyVoidType(t *irtypes.VoidType) *irtypes.VoidType {
	return &irtypes.VoidType{
		TypeName: t.TypeName,
	}
}

// copyFuncType returns an identical copy of the given LLVM IR type.
func copyFuncType(t *irtypes.FuncType) *irtypes.FuncType {
	var params []irtypes.Type
	for i := range t.Params {
		param := copyType(t.Params[i])
		params = append(params, param)
	}
	return &irtypes.FuncType{
		TypeName: t.TypeName,
		RetType:  copyType(t.RetType),
		Params:   params,
		Variadic: t.Variadic,
	}
}

// copyIntType returns an identical copy of the given LLVM IR type.
func copyIntType(t *irtypes.IntType) *irtypes.IntType {
	return &irtypes.IntType{
		TypeName: t.TypeName,
		BitSize:  t.BitSize,
	}
}

// copyFloatType returns an identical copy of the given LLVM IR type.
func copyFloatType(t *irtypes.FloatType) *irtypes.FloatType {
	return &irtypes.FloatType{
		TypeName: t.TypeName,
		Kind:     t.Kind,
	}
}

// copyMMXType returns an identical copy of the given LLVM IR type.
func copyMMXType(t *irtypes.MMXType) *irtypes.MMXType {
	return &irtypes.MMXType{
		TypeName: t.TypeName,
	}
}

// copyPointerType returns an identical copy of the given LLVM IR type.
func copyPointerType(t *irtypes.PointerType) *irtypes.PointerType {
	return &irtypes.PointerType{
		TypeName:  t.TypeName,
		ElemType:  copyType(t.ElemType),
		AddrSpace: t.AddrSpace,
	}
}

// copyVectorType returns an identical copy of the given LLVM IR type.
func copyVectorType(t *irtypes.VectorType) *irtypes.VectorType {
	return &irtypes.VectorType{
		TypeName: t.TypeName,
		Scalable: t.Scalable,
		Len:      t.Len,
		ElemType: copyType(t.ElemType),
	}
}

// copyLabelType returns an identical copy of the given LLVM IR type.
func copyLabelType(t *irtypes.LabelType) *irtypes.LabelType {
	return &irtypes.LabelType{
		TypeName: t.TypeName,
	}
}

// copyTokenType returns an identical copy of the given LLVM IR type.
func copyTokenType(t *irtypes.TokenType) *irtypes.TokenType {
	return &irtypes.TokenType{
		TypeName: t.TypeName,
	}
}

// copyMetadataType returns an identical copy of the given LLVM IR type.
func copyMetadataType(t *irtypes.MetadataType) *irtypes.MetadataType {
	return &irtypes.MetadataType{
		TypeName: t.TypeName,
	}
}

// copyArrayType returns an identical copy of the given LLVM IR type.
func copyArrayType(t *irtypes.ArrayType) *irtypes.ArrayType {
	return &irtypes.ArrayType{
		TypeName: t.TypeName,
		Len:      t.Len,
		ElemType: copyType(t.ElemType),
	}
}

// copyStructType returns an identical copy of the given LLVM IR type.
func copyStructType(t *irtypes.StructType) *irtypes.StructType {
	var fields []irtypes.Type
	for i := range t.Fields {
		field := copyType(t.Fields[i])
		fields = append(fields, field)
	}
	return &irtypes.StructType{
		TypeName: t.TypeName,
		Packed:   t.Packed,
		Fields:   fields,
		Opaque:   t.Opaque,
	}
}

// ### [ Helper functions ] ####################################################

// newSliceType returns a new LLVM IR slice type based on the given element
// type.
func (m *Module) newSliceType(elemType irtypes.Type) *irtypes.StructType {
	typeName := m.getSliceTypeName(elemType)
	if typ, ok := m.types[typeName]; ok {
		return typ.(*irtypes.StructType)
	}
	data := irtypes.NewPointer(elemType)
	length := m.irTypeFromName("int")
	capacity := m.irTypeFromName("int")
	// TODO: add support for LLVM IR structure types with field names.
	//sliceType = NewStruct(
	//   Field{Name: "data", Type: data},
	//   Field{Name: "len", Type: length},
	//   Field{Name: "cap", Type: capacity},
	//)
	typ := irtypes.NewStruct(data, length, capacity)
	typ.SetName(typeName)
	m.types[typeName] = typ
	m.Module.TypeDefs = append(m.Module.TypeDefs, typ)
	return typ
}

// getSliceTypeName returns the LLVM IR type name of the slice type with the
// specified element type.
func (m *Module) getSliceTypeName(elemType irtypes.Type) string {
	return "[]" + elemType.String() // TODO: use fully qualified name for elemType.
}
