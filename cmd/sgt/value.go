package main

import (
	"fmt"
	goconstant "go/constant"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"golang.org/x/tools/go/ssa"
)

// TODO: add defValue analogous to useValue?

// defValue stores the source value to the given Go destination, emitting a
// store instruction to fn.
func (gen *generator) defValue(goDst ssa.Value, src value.Value, fn *Func) {
	dst := gen.llValueFromGoValue(goDst)
	switch dst := dst.(type) {
	// Global variable (with type pointer to content type).
	case *ir.Global:
		storeInst := fn.Cur.NewStore(src, dst)
		dbg.Println("store instruction:", storeInst.LLString())
	// Local variable (with type pointer to content type).
	//case *ir.InstAlloca:
	//	storeInst := fn.Cur.NewStore(src, dst)
	//	dbg.Println("store instruction:", storeInst.LLString())
	default:
		panic(fmt.Errorf("support for Go dst value %T not yet implemented", dst))
	}
}

// useValue returns the LLVM IR value corresponding to the given Go value. If
// the LLVM IR value is a global variable (with type pointer to content) or a
// local variable (alloca, with type pointer to concent), the load the content
// value before return, emitting a load instruction to fn.
func (gen *generator) useValue(goValue ssa.Value, fn *Func) value.Value {
	switch goValue := goValue.(type) {
	// Local variables.
	case *ssa.UnOp:
		local := fn.locals[goValue.Name()]
		//pretty.Println("fn.locals:", fn.locals)
		dbg.Println("local:", local.LLString())
		return fn.Cur.NewLoad(local.Type().(*types.PointerType).ElemType, local)
	}
	v := gen.llValueFromGoValue(goValue)
	switch v := v.(type) {
	// Global variable (with type pointer to content type).
	case *ir.Global:
		loadInst := fn.Cur.NewLoad(v.ContentType, v)
		return loadInst
	// Local variable (with type pointer to content type).
	//case *ir.InstAlloca:
	//	return fn.Cur.NewLoad(v.Typ.ElemType, v)
	default:
		return v
	}
}

// llValueFromGoValue returns the LLVM IR value corresponding to the given Go
// value.
func (gen *generator) llValueFromGoValue(goValue ssa.Value) value.Value {
	switch goValue := goValue.(type) {
	case *ssa.Const:
		return gen.llValueFromGoConst(goValue)
	case *ssa.Global:
		return gen.llValueFromGoGlobal(goValue)
	default:
		panic(fmt.Errorf("support for Go value %T not yet implemented", goValue))
	}
}

// --- [ constant ] ------------------------------------------------------------

// llValueFromGoConst returns the LLVM IR constant corresponding to the given Go
// constant.
func (gen *generator) llValueFromGoConst(goValue *ssa.Const) irconstant.Constant {
	typ := gen.llTypeFromGoType(goValue.Type())
	_ = typ // TODO: figure out how to handle type.
	goVal := goconstant.Val(goValue.Value)
	switch goVal := goVal.(type) {
	case bool:
		return irconstant.NewBool(goVal)
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
		zero := irconstant.NewInt(types.I64, 0)
		data := irconstant.NewGetElementPtr(strLit.Typ, g, zero, zero)
		// Create len field.
		length := irconstant.NewInt(lenType, n)
		// Return LLVM IR string constant.
		return irconstant.NewStruct(stringType, data, length)
	default:
		panic(fmt.Errorf("support for converting Go string literal to LLVM IR constant with LLVM IR type %T (%q) not yet implemented", typ, typ.Name()))
	}
}

// --- [ global ] --------------------------------------------------------------

// llValueFromGoGlobal returns the LLVM IR global variable corresponding to the
// given Go global.
func (gen *generator) llValueFromGoGlobal(goValue *ssa.Global) *ir.Global {
	global, ok := gen.globals[goValue.Name()]
	if !ok {
		panic(fmt.Errorf("unable to locate LLVM IR global variable corresponding to Go global with name %q", goValue.Name()))
	}
	return global
}
