package main

import (
	"fmt"
	goconstant "go/constant"

	"github.com/llir/llvm/ir/constant"
	irconstant "github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"golang.org/x/tools/go/ssa"
)

// llValueFromGoValue returns the LLVM IR value corresponding to the given Go
// value.
func (gen *generator) llValueFromGoValue(goValue ssa.Value) value.Value {
	switch goValue := goValue.(type) {
	case *ssa.Const:
		return gen.llValueFromGoConst(goValue)
	default:
		panic(fmt.Errorf("support for Go value %T not yet implemented", goValue))
	}
}

// llValueFromGoConst returns the LLVM IR constant corresponding to the given Go
// constant.
func (gen *generator) llValueFromGoConst(goValue *ssa.Const) irconstant.Constant {
	typ := gen.llTypeFromGoType(goValue.Type())
	_ = typ // TODO: figure out how to handle type.
	goVal := goconstant.Val(goValue.Value)
	switch goVal := goVal.(type) {
	case string:
		return gen.llValueFromGoStringLit(typ, goVal)
	case nil:
		// check constant kind for nil values, as go/constant.Val returns nil also
		// for complex literals.
		break
	default:
		panic(fmt.Errorf("support for Go constant literal %T not yet implemented", goVal))
	}
	switch kind := goValue.Value.Kind(); kind {
	default:
		panic(fmt.Errorf("support for Go constant kind %v not yet implemented", kind))
	}
}

// llValueFromGoStringLit returns the LLVM IR constant corresponding to the
// given Go string literal.
func (gen *generator) llValueFromGoStringLit(typ irtypes.Type, s string) irconstant.Constant {
	switch typ.Name() {
	case "string":
		g := gen.globalFromStringLit(s)
		strLit := g.Init.(*irconstant.CharArray)
		n := int64(len(strLit.X))
		// Unpack %string type.
		stringType := typ.(*irtypes.StructType)
		dataType := stringType.Fields[0].(*irtypes.PointerType)
		_ = dataType
		lenType := stringType.Fields[1].(*irtypes.IntType)
		// Create data field.
		zero := constant.NewInt(types.I64, 0)
		data := constant.NewGetElementPtr(strLit.Typ, g, zero, zero)
		// Create len field.
		length := constant.NewInt(lenType, n)
		// Return LLVM IR string constant.
		return irconstant.NewStruct(stringType, data, length)
	default:
		panic(fmt.Errorf("support for converting Go string literal to LLVM IR constant with LLVM IR type %T (%q) not yet implemented", typ, typ.Name()))
	}
}
