package irgen

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
)

// synthLen synthesizes a builtin `len` function based on the given argument
// type, emitting to m.
func (m *Module) synthLen(argType irtypes.Type) *ir.Func {
	dbg.Println("synthLen")
	// Define `len(T)` function if not present.
	typeName := argType.Name()
	lenFuncName := fmt.Sprintf("len(%s)", typeName)
	if lenFunc, ok := m.predeclaredFuncs[lenFuncName]; ok {
		return lenFunc
	}
	retType := m.irTypeFromName("int")
	arg := ir.NewParam("v", argType)
	lenFunc := m.Module.NewFunc(lenFuncName, retType, arg)
	entry := lenFunc.NewBlock("entry")
	var length irvalue.Value
	switch argType := argType.(type) {
	case *irtypes.StructType:
		switch {
		// string
		case argType.Name() == "string":
			lengthField := entry.NewExtractValue(arg, 1)
			addMetadata(lengthField, "field", "len")
			length = lengthField
		// slice
		case strings.HasPrefix(argType.Name(), "[]"):
			lengthField := entry.NewExtractValue(arg, 1)
			addMetadata(lengthField, "field", "len")
			length = lengthField
		default:
			panic(fmt.Errorf("support for type %T (%q) as argument to builtin len function not yet implemented", argType, argType.Name()))
		}
	default:
		panic(fmt.Errorf("support for type %T (%q) as argument to builtin len function not yet implemented", argType, argType.Name()))
	}
	entry.NewRet(length)
	m.predeclaredFuncs[lenFuncName] = lenFunc
	return lenFunc
}
