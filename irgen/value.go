package irgen

import (
	"fmt"
	goconstant "go/constant"
	"math/big"
	"strconv"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
	"golang.org/x/tools/go/ssa"
)

// --- [ use ] -----------------------------------------------------------------

// useValue returns the LLVM IR value corresponding to the given local or global
// Go SSA value, emitting to fn and fn.m. If the LLVM IR value is a global
// variable (with pointer to content type) then load the content value before
// return, emitting a load instruction to fn.
//
// Pre-condition: index global and local values of fn and fn.m.
func (fn *Func) useValue(goValue ssa.Value) irvalue.Value {
	v := fn.irValueFromGo(goValue)
	if global, ok := v.(*ir.Global); ok {
		// Load contents of global variable for use.
		return fn.cur.NewLoad(global.ContentType, global)
	}
	return v
}

// --- [ convert ] -------------------------------------------------------------

// irValueFromGo returns the LLVM IR value corresponding to the given local or
// global Go SSA value, emitting to fn and fn.m.
//
// Pre-condition: index global and local values of fn and fn.m.
func (fn *Func) irValueFromGo(goValue ssa.Value) irvalue.Value {
	// Translate local or global Go SSA value.
	switch goValue := goValue.(type) {
	case *ssa.FreeVar:
		panic("support for *ssa.FreeVar not yet implemented")
	case *ssa.Parameter:
		// Lookup indexed LLVM IR function parameter of Go SSA function parameter.
		if v, ok := fn.locals[goValue]; ok {
			return v
		}
		// Pre-condition invalidated, function parameter should have been indexed.
		// This is a fatal error and indicates a bug in irgen.
		panic(fmt.Errorf("unable to locate indexed LLVM IR function parameter of Go SSA function parameter %q", goValue.Name()))
	// Value instruction.
	case ssaValueInstruction:
		// Lookup indexed LLVM IR value of Go SSA value instruction.
		if v, ok := fn.locals[goValue]; ok {
			return v
		}
		// Pre-condition invalidated, value instruction should have been indexed.
		// This is a fatal error and indicates a bug in irgen.
		panic(fmt.Errorf("unable to locate indexed LLVM IR value of Go SSA value instruction %q", goValue.Name()))
	default:
		return fn.m.irValueFromGo(goValue)
	}
}

// irValueFromGo returns the LLVM IR value corresponding to the given global Go
// SSA value, emitting to m.
func (m *Module) irValueFromGo(goValue ssa.Value) irvalue.Value {
	// Lookup indexed LLVM IR value of global Go SSA value.
	if v, ok := m.globals[goValue]; ok {
		return v
	}
	switch goValue := goValue.(type) {
	// Global Go SSA values.
	case *ssa.Builtin:
		return m.irValueFromGoBuiltin(goValue)
	case *ssa.Const:
		return m.irValueFromGoConst(goValue)
	case *ssa.Function:
		return m.irValueFromGoFunc(goValue)
	case *ssa.Global:
		warn.Printf("unable to locate LLVM IR value of global Go value %q", m.fullName(goValue))
		panic("support for *ssa.Global not yet implemented")
	default:
		panic(fmt.Errorf("support for Go SSA value %T not yet implemented", goValue))
	}
}

// --- [ builtin ] -------------------------------------------------------------

// irValueFromGoBuiltin returns the LLVM IR value corresponding to the given Go
// SSA builtin value, emitting to m.
func (m *Module) irValueFromGoBuiltin(goValue *ssa.Builtin) irvalue.Value {
	dbg.Println("irValueFromGoBuiltin")
	f, ok := m.predeclaredFuncs[goValue.Name()]
	if !ok {
		panic(fmt.Errorf("unable to locate LLVM IR value of Go builtin value %q", goValue.Name()))
	}
	return f
}

// --- [ constant ] ------------------------------------------------------------

// irValueFromGoConst returns the LLVM IR constant corresponding to the given Go
// SSA constant, emitting to m.
func (m *Module) irValueFromGoConst(goConst *ssa.Const) irconstant.Constant {
	dbg.Println("irValueFromGoConst")
	typ := m.irTypeFromGo(goConst.Type())
	dbg.Println("   typ:", typ)
	goVal := goconstant.Val(goConst.Value)
	dbg.Println("   goVal:", goVal)
	switch goVal := goVal.(type) {
	// boolean literal
	case bool:
		return irconstant.NewBool(goVal)
	// string literal
	case string:
		return m.irValueFromGoStringLit(typ, goVal)
	// integer literal
	case int64:
		switch typ := typ.(type) {
		case *irtypes.IntType:
			return irconstant.NewInt(typ, goVal)
		case *irtypes.FloatType:
			s := strconv.FormatInt(goVal, 10) + ".0"
			c, err := irconstant.NewFloatFromString(typ, s)
			if err != nil {
				panic(fmt.Errorf("unable to parse floating-point literal %q", s))
			}
			return c
		default:
			panic(fmt.Errorf("support for integer literal of type %T not yet implemented", typ))
		}
	case *big.Int:
		x := big.NewInt(0).Set(goVal)
		return &irconstant.Int{
			Typ: typ.(*irtypes.IntType),
			X:   x,
		}
	// floating-point literal
	case *big.Rat:
		floatType := typ.(*irtypes.FloatType)
		prec := precFromFloatKind(floatType.Kind)
		x := big.NewFloat(0).SetPrec(prec)
		x.SetRat(goVal)
		return &irconstant.Float{
			Typ: floatType,
			X:   x,
		}
	case *big.Float:
		floatType := typ.(*irtypes.FloatType)
		prec := precFromFloatKind(floatType.Kind)
		x := big.NewFloat(0).SetPrec(prec)
		x.Set(goVal)
		return &irconstant.Float{
			Typ: floatType,
			X:   x,
		}
	// kind of everything else is nil
	case nil:
		// Check constant kind for nil values, as go/constant.Val returns nil also
		// for complex constant literals.
		switch kind := goConst.Value.Kind(); kind {
		default:
			panic(fmt.Errorf("support for Go constant kind %v not yet implemented", kind))
		}
	default:
		panic(fmt.Errorf("support for Go constant %T not yet implemented", goVal))
	}
}

// ~~~ [ string literal ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// irValueFromGoStringLit returns the LLVM IR constant corresponding to the
// given Go string literal, emitting to m.
func (m *Module) irValueFromGoStringLit(typ irtypes.Type, s string) irconstant.Constant {
	switch typ.Name() {
	case "string":
		g := m.emitStringLit(s)
		strLit := g.Init.(*irconstant.CharArray)
		n := int64(len(strLit.X))
		// Unpack %string type.
		stringType := typ.(*irtypes.StructType)
		dataType := stringType.Fields[0].(*irtypes.PointerType)
		_ = dataType
		lenType := stringType.Fields[1].(*irtypes.IntType)
		// Create data field.
		zero := irconstant.NewInt(irtypes.I64, 0)
		data := irconstant.NewGetElementPtr(strLit.Typ, g, zero, zero)
		// Create len field.
		length := irconstant.NewInt(lenType, n)
		// Return LLVM IR string constant.
		return irconstant.NewStruct(stringType, data, length)
	default:
		panic(fmt.Errorf("support for converting Go string literal to LLVM IR constant with LLVM IR type %T (%q) not yet implemented", typ, typ.Name()))
	}
}

// --- [ function ] ------------------------------------------------------------

// irValueFromGoFunc returns the LLVM IR function corresponding to the given Go
// SSA function, emitting to m.
func (m *Module) irValueFromGoFunc(goFunc *ssa.Function) *ir.Func {
	dbg.Println("irValueFromGoFunc")
	return m.getFunc(goFunc)
}
