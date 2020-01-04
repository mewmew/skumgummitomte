package irgen

import (
	"fmt"
	"go/token"
	gotypes "go/types"
	"strings"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	irenum "github.com/llir/llvm/ir/enum"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// irValueInstruction is an LLVM IR value instruction.
type irValueInstruction interface {
	irvalue.Named
	ir.Instruction
}

// ssaValueInstruction is a Go SSA value instruction.
type ssaValueInstruction interface {
	ssa.Value
	ssa.Instruction
}

// emitInst compiles the given Go SSA instruction to corresponding LLVM IR
// instructions, emitting to fn.
func (fn *Func) emitInst(goInst ssa.Instruction) error {
	switch goInst := goInst.(type) {
	// Value-producing instructions (local Go SSA values).
	case ssaValueInstruction:
		return fn.emitValueInst(goInst)
	// Non-value producing instructions.
	case *ssa.DebugRef:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic("support for *ssa.DebugRef not yet implemented")
	case *ssa.Defer:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic("support for *ssa.Defer not yet implemented")
	case *ssa.Go:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic("support for *ssa.Go not yet implemented")
	case *ssa.If:
		return fn.emitIf(goInst)
	case *ssa.Jump:
		return fn.emitJump(goInst)
	case *ssa.MapUpdate:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic("support for *ssa.MapUpdate not yet implemented")
	case *ssa.Panic:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic("support for *ssa.Panic not yet implemented")
	case *ssa.Return:
		return fn.emitReturn(goInst)
	case *ssa.RunDefers:
		// TODO: implement support for defer.
		return nil // ignore *ssa.RunDefers instruction for now
	case *ssa.Send:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic("support for *ssa.Send not yet implemented")
	case *ssa.Store:
		return fn.emitStore(goInst)
	default:
		panic(fmt.Errorf("support for Go SSA instruction %T not yet implemented", goInst))
	}
}

// emitValueInst compiles the given Go SSA value instruction to corresponding
// LLVM IR instructions, emitting to fn.
func (fn *Func) emitValueInst(goInst ssaValueInstruction) error {
	switch goInst := goInst.(type) {
	case *ssa.Alloc:
		return fn.emitAlloc(goInst)
	case *ssa.BinOp:
		return fn.emitBinOp(goInst)
	case *ssa.Call:
		return fn.emitCall(goInst)
	case *ssa.ChangeInterface:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.ChangeInterface (in %q) not yet implemented", goInst.Name()))
	case *ssa.ChangeType:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.ChangeType (in %q) not yet implemented", goInst.Name()))
	case *ssa.Convert:
		return fn.emitConvert(goInst)
	case *ssa.Extract:
		return fn.emitExtract(goInst)
	case *ssa.Field:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.Field (in %q) not yet implemented", goInst.Name()))
	case *ssa.FieldAddr:
		return fn.emitFieldAddr(goInst)
	case *ssa.Index:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.Index (in %q) not yet implemented", goInst.Name()))
	case *ssa.IndexAddr:
		return fn.emitIndexAddr(goInst)
	case *ssa.Lookup:
		return fn.emitLookup(goInst)
	case *ssa.MakeChan:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.MakeChan (in %q) not yet implemented", goInst.Name()))
	case *ssa.MakeClosure:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.MakeClosure (in %q) not yet implemented", goInst.Name()))
	case *ssa.MakeInterface:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.MakeInterface (in %q) not yet implemented", goInst.Name()))
	case *ssa.MakeMap:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.MakeMap (in %q) not yet implemented", goInst.Name()))
	case *ssa.MakeSlice:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.MakeSlice (in %q) not yet implemented", goInst.Name()))
	case *ssa.Next:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.Next (in %q) not yet implemented", goInst.Name()))
	case *ssa.Phi:
		return fn.emitPhi(goInst)
	case *ssa.Range:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.Range (in %q) not yet implemented", goInst.Name()))
	case *ssa.Select:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.Select (in %q) not yet implemented", goInst.Name()))
	case *ssa.Slice:
		return fn.emitSlice(goInst)
	case *ssa.TypeAssert:
		goInst.Parent().WriteTo(ssaDebugWriter)
		panic(fmt.Errorf("support for *ssa.TypeAssert (in %q) not yet implemented", goInst.Name()))
	case *ssa.UnOp:
		return fn.emitUnOp(goInst)
	default:
		panic(fmt.Errorf("support for Go SSA value instruction %T not yet implemented", goInst))
	}
}

// === [ Non-value instructions ] ==============================================

// --- [ if instruction ] ------------------------------------------------------

// emitIf compiles the given Go SSA if instruction to corresponding LLVM IR
// instructions, emitting to fn.
func (fn *Func) emitIf(goInst *ssa.If) error {
	dbg.Println("emitIf")
	cond := fn.useValue(goInst.Cond)
	dbg.Println("   cond:", cond)
	// The If instruction transfers control to one of the two successors of its
	// owning block, depending on the boolean Cond: the first if true, the second
	// if false.
	succs := goInst.Block().Succs
	targetTrue := fn.getBlock(succs[0])
	targetFalse := fn.getBlock(succs[1])
	term := fn.cur.NewCondBr(cond, targetTrue, targetFalse)
	dbg.Println("   term:", term.LLString())
	return nil
}

// --- [ jump instruction ] ----------------------------------------------------

// emitJump compiles the given Go SSA jump instruction to corresponding LLVM IR
// instructions, emitting to fn.
func (fn *Func) emitJump(goInst *ssa.Jump) error {
	dbg.Println("emitJump")
	// The Jump instruction transfers control to the sole successor of its owning block.
	succs := goInst.Block().Succs
	target := fn.getBlock(succs[0])
	term := fn.cur.NewBr(target)
	dbg.Println("   term:", term.LLString())
	return nil
}

// --- [ return instruction ] --------------------------------------------------

// emitReturn compiles the given Go SSA return instruction to corresponding LLVM
// IR instructions, emitting to fn.
func (fn *Func) emitReturn(goInst *ssa.Return) error {
	dbg.Println("emitReturn")
	// Results of function.
	var results []irvalue.Value
	for _, goResult := range goInst.Results {
		result := fn.useValue(goResult)
		dbg.Println("   result:", result)
		if r, ok := result.(*irconstant.Float); ok {
			dbg.Println("   result (float):", r.X.String())
		}
		results = append(results, result)
	}
	// Return value; nil if void return, and struct if multiple return value.
	var x irvalue.Value
	switch len(results) {
	// void return.
	case 0:
		x = nil // void return
	// single value return.
	case 1:
		x = results[0]
	// multiple value return.
	default:
		structType, ok := fn.Func.Sig.RetType.(*irtypes.StructType)
		if !ok {
			return errors.Errorf("invalid return type for function %q with multiple return values (%d); expected *irtypes.StructType, got %T", fn.Func.Name(), len(results), fn.Func.Sig.RetType)
		}
		if len(structType.Fields) != len(results) {
			return errors.Errorf("mismatch between number of results in function signature (%d) and function return values (%d) in function %q", len(structType.Fields), len(results))
		}
		alloca := fn.entry.NewAlloca(structType)
		fn.cur.NewStore(irconstant.NewZeroInitializer(structType), alloca)
		var ret irValueInstruction = fn.cur.NewLoad(structType, alloca)
		for index, result := range results {
			ret = fn.cur.NewInsertValue(ret, result, uint64(index))
		}
		x = ret
	}
	term := fn.cur.NewRet(x)
	dbg.Println("   term:", term.LLString())
	return nil
}

// --- [ store instruction ] ---------------------------------------------------

// emitStore compiles the given Go SSA store instruction to corresponding LLVM
// IR instructions, emitting to fn.
func (fn *Func) emitStore(goInst *ssa.Store) error {
	dbg.Println("emitStore")
	addr := fn.useValue(goInst.Addr)
	dbg.Println("   addr:", addr)
	val := fn.useValue(goInst.Val)
	dbg.Println("   val:", val)
	inst := fn.cur.NewStore(val, addr)
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// === [ Value instructions ] ==================================================

// --- [ alloc instruction ] ---------------------------------------------------

// emitAlloc compiles the given Go SSA alloc instruction to corresponding LLVM
// IR instructions, emitting to fn.
func (fn *Func) emitAlloc(goInst *ssa.Alloc) error {
	dbg.Println("emitAlloc")
	if goInst.Heap {
		// Allocate space in the heap.
		return fn.emitNew(goInst)
	}
	typ := fn.m.irTypeFromGo(goInst.Type())
	ptrType := typ.(*irtypes.PointerType)
	inst := fn.cur.NewAlloca(ptrType.ElemType)
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	// Add local variable name metadata attachment to alloca instruction.
	if len(goInst.Comment) > 0 {
		addMetadata(inst, "var_name", goInst.Comment)
	}
	dbg.Println("   inst:", inst.LLString())
	// The space of a stack allocated local variable is re-initialized to zero
	// each time it is executed.
	zero := irconstant.NewZeroInitializer(ptrType.ElemType)
	fn.cur.NewStore(zero, inst)
	return nil
}

// ~~~ [ new - heap alloc instruction ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TODO: move emitNew to builtin, where other synthesized builtin functions are
// handled.

// emitNew compiles the given Go SSA heap alloc instruction to corresponding
// LLVM IR instructions, emitting to fn.
func (fn *Func) emitNew(goInst *ssa.Alloc) error {
	dbg.Println("emitNew")
	typ := fn.m.irTypeFromGo(goInst.Type())
	ptrType := typ.(*irtypes.PointerType)
	// Define `new(T)` function if not present.
	typeName := goInst.Type().(*gotypes.Pointer).Elem().String()
	newFuncName := fmt.Sprintf("new(%s)", typeName)
	newFunc, ok := fn.m.predeclaredFuncs[newFuncName]
	if !ok {
		retType := ptrType
		newFunc = fn.m.Module.NewFunc(newFuncName, retType)
		entry := newFunc.NewBlock("entry")
		allocaInst := entry.NewAlloca(ptrType.ElemType)
		objectsizeFunc := fn.m.getPredeclaredFunc("llvm.objectsize.i64")
		bitCastInst := entry.NewBitCast(allocaInst, irtypes.I8Ptr)
		objectsizeArgs := []irvalue.Value{
			bitCastInst,      // object
			irconstant.False, // min
			irconstant.False, // nullunknown
			irconstant.False, // dynamic
		}
		size := entry.NewCall(objectsizeFunc, objectsizeArgs...)
		size.SetName("size")
		cond := entry.NewICmp(irenum.IPredNE, size, irconstant.NewInt(irtypes.I64, -1))
		success := newFunc.NewBlock("success")
		fail := newFunc.NewBlock("fail")
		entry.NewCondBr(cond, success, fail)
		// Generate `success` basic block.
		callocFunc := fn.m.getPredeclaredFunc("calloc") // using calloc to zero initialize
		args := []irvalue.Value{
			irconstant.NewInt(irtypes.I64, 1),
			size,
		}
		callInst := success.NewCall(callocFunc, args...)
		result := success.NewBitCast(callInst, ptrType)
		success.NewRet(result)
		// Generate `fail` basic block.
		// TODO: panic with "unable to get size of type T" error message.
		fail.NewUnreachable()
		// Add synthesized `new(T)` function to predeclared functions.
		fn.m.predeclaredFuncs[newFunc.Name()] = newFunc
	}
	// Invoke new(T).
	inst := fn.cur.NewCall(newFunc)
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	// Add local variable name metadata attachment to alloca instruction.
	if len(goInst.Comment) > 0 {
		addMetadata(inst, "comment", goInst.Comment)
	}
	dbg.Println("   inst:", inst)
	return nil
}

// --- [ binary operation instruction ] ----------------------------------------

// emitBinOp compiles the given Go SSA binary operation instruction to
// corresponding LLVM IR instructions, emitting to fn.
func (fn *Func) emitBinOp(goInst *ssa.BinOp) error {
	dbg.Println("emitBinOp")
	dbg.Println("   op:", goInst.Op)
	x := fn.useValue(goInst.X)
	dbg.Println("   x:", x)
	y := fn.useValue(goInst.Y)
	dbg.Println("   y:", y)
	var inst irValueInstruction
	switch goInst.Op {
	// ADD (+)
	case token.ADD: // +
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			inst = fn.cur.NewAdd(x, y)
		case *irtypes.FloatType:
			inst = fn.cur.NewFAdd(x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "complex64", "complex128":
				op := func(a, b irvalue.Value) irValueInstruction {
					return fn.cur.NewFAdd(a, b)
				}
				inst = fn.emitComplexBinOp(op, x, y)
			case "string":
				// TODO: add support for string type.
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// SUB (-)
	case token.SUB: // -
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			inst = fn.cur.NewSub(x, y)
		case *irtypes.FloatType:
			inst = fn.cur.NewFSub(x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "complex64", "complex128":
				op := func(a, b irvalue.Value) irValueInstruction {
					return fn.cur.NewFSub(a, b)
				}
				inst = fn.emitComplexBinOp(op, x, y)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// MUL (*)
	case token.MUL: // *
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			inst = fn.cur.NewMul(x, y)
		case *irtypes.FloatType:
			inst = fn.cur.NewFMul(x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "complex64", "complex128":
				op := func(a, b irvalue.Value) irValueInstruction {
					return fn.cur.NewFMul(a, b)
				}
				inst = fn.emitComplexBinOp(op, x, y)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// QUO (/)
	case token.QUO: // /
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewSDiv(x, y)
			} else {
				inst = fn.cur.NewUDiv(x, y)
			}
		case *irtypes.FloatType:
			inst = fn.cur.NewFDiv(x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "complex64", "complex128":
				op := func(a, b irvalue.Value) irValueInstruction {
					return fn.cur.NewFDiv(a, b)
				}
				inst = fn.emitComplexBinOp(op, x, y)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// REM (%)
	case token.REM: // %
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewSRem(x, y)
			} else {
				inst = fn.cur.NewURem(x, y)
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// AND (&)
	case token.AND: // &
		inst = fn.cur.NewAnd(x, y)
	// OR (|)
	case token.OR: // |
		inst = fn.cur.NewOr(x, y)
	// XOR (^)
	case token.XOR: // ^
		inst = fn.cur.NewXor(x, y)
	// SHL (<<)
	case token.SHL: // <<
		inst = fn.cur.NewShl(x, y)
	// SHR (>>)
	case token.SHR: // >>
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewAShr(x, y)
			} else {
				inst = fn.cur.NewLShr(x, y)
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// AND_NOT (&^)
	case token.AND_NOT: // &^
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			zero := irconstant.NewInt(typ, 0)
			tmp := fn.cur.NewXor(y, zero)
			dbg.Println("   tmp:", tmp.LLString())
			inst = fn.cur.NewAnd(x, tmp)
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// EQL (==)
	case token.EQL: // ==
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			inst = fn.cur.NewICmp(irenum.IPredEQ, x, y)
		case *irtypes.FloatType:
			// TODO: figure out when to use FPredOEQ vs. FPredUEQ (ordered vs.
			// unordered).
			inst = fn.cur.NewFCmp(irenum.FPredOEQ, x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "complex64", "complex128":
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			case "string":
				cmp := fn.m.getPredeclaredFunc("cmp.string")
				result := fn.cur.NewCall(cmp, x, y)
				zero := irconstant.NewInt(result.Type().(*irtypes.IntType), 0)
				inst = fn.cur.NewICmp(irenum.IPredEQ, result, zero)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// NEQ (!=)
	case token.NEQ: // !=
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			inst = fn.cur.NewICmp(irenum.IPredNE, x, y)
		case *irtypes.FloatType:
			// TODO: figure out when to use FPredONE vs. FPredUNE (ordered vs.
			// unordered).
			inst = fn.cur.NewFCmp(irenum.FPredONE, x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "complex64", "complex128":
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			case "string":
				cmp := fn.m.getPredeclaredFunc("cmp.string")
				result := fn.cur.NewCall(cmp, x, y)
				zero := irconstant.NewInt(result.Type().(*irtypes.IntType), 0)
				inst = fn.cur.NewICmp(irenum.IPredNE, result, zero)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// LSS (<)
	case token.LSS: // <
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewICmp(irenum.IPredSLT, x, y)
			} else {
				inst = fn.cur.NewICmp(irenum.IPredULT, x, y)
			}
		case *irtypes.FloatType:
			// TODO: figure out when to use FPredOLT vs. FPredULT (ordered vs.
			// unordered).
			inst = fn.cur.NewFCmp(irenum.FPredOLT, x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "string":
				cmp := fn.m.getPredeclaredFunc("cmp.string")
				result := fn.cur.NewCall(cmp, x, y)
				zero := irconstant.NewInt(result.Type().(*irtypes.IntType), -1)
				inst = fn.cur.NewICmp(irenum.IPredEQ, result, zero)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// LEQ (<=)
	case token.LEQ: // <=
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewICmp(irenum.IPredSLE, x, y)
			} else {
				inst = fn.cur.NewICmp(irenum.IPredULE, x, y)
			}
		case *irtypes.FloatType:
			// TODO: figure out when to use FPredOLE vs. FPredULE (ordered vs.
			// unordered).
			inst = fn.cur.NewFCmp(irenum.FPredOLE, x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "string":
				cmp := fn.m.getPredeclaredFunc("cmp.string")
				result := fn.cur.NewCall(cmp, x, y)
				zero := irconstant.NewInt(result.Type().(*irtypes.IntType), 1)
				inst = fn.cur.NewICmp(irenum.IPredNE, result, zero)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// GTR (>)
	case token.GTR: // >
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewICmp(irenum.IPredSGT, x, y)
			} else {
				inst = fn.cur.NewICmp(irenum.IPredUGT, x, y)
			}
		case *irtypes.FloatType:
			// TODO: figure out when to use FPredOGT vs. FPredUGT (ordered vs.
			// unordered).
			inst = fn.cur.NewFCmp(irenum.FPredOGT, x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "string":
				cmp := fn.m.getPredeclaredFunc("cmp.string")
				result := fn.cur.NewCall(cmp, x, y)
				zero := irconstant.NewInt(result.Type().(*irtypes.IntType), 1)
				inst = fn.cur.NewICmp(irenum.IPredEQ, result, zero)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// GEQ (>=)
	case token.GEQ: // >=
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			if fn.m.isSigned(typ) {
				inst = fn.cur.NewICmp(irenum.IPredSGE, x, y)
			} else {
				inst = fn.cur.NewICmp(irenum.IPredUGE, x, y)
			}
		case *irtypes.FloatType:
			// TODO: figure out when to use FPredOGE vs. FPredUGE (ordered vs.
			// unordered).
			inst = fn.cur.NewFCmp(irenum.FPredOGE, x, y)
		case *irtypes.StructType:
			switch typ.Name() {
			case "string":
				cmp := fn.m.getPredeclaredFunc("cmp.string")
				result := fn.cur.NewCall(cmp, x, y)
				zero := irconstant.NewInt(result.Type().(*irtypes.IntType), -1)
				inst = fn.cur.NewICmp(irenum.IPredNE, result, zero)
			default:
				panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
			}
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	default:
		panic(fmt.Errorf("support for Go SSA binary operation instruction token %v not yet implemented", goInst.Op))
	}
	inst.SetName(goInst.Name())
	// Add binary operation token metadata attachment to instruction.
	addMetadata(inst, "binary_op", goInst.Op.String())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// ~~~ [ complex binary operation ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// emitComplexBinOp compiles the given complex binary operation to corresponding
// LLVM IR instructions, emitting to fn.
func (fn *Func) emitComplexBinOp(op func(a, b irvalue.Value) irValueInstruction, x, y irvalue.Value) irValueInstruction {
	typ := x.Type()
	// real part.
	xReal := fn.cur.NewExtractValue(x, 0)
	addMetadata(xReal, "comment", "real")
	dbg.Println("   xReal:", xReal.LLString())
	yReal := fn.cur.NewExtractValue(y, 0)
	addMetadata(yReal, "comment", "real")
	dbg.Println("   yReal:", yReal.LLString())
	real := op(xReal, yReal)
	addMetadata(real, "comment", "real")
	dbg.Println("   real:", real.LLString())
	// imaginary part.
	xImag := fn.cur.NewExtractValue(x, 1)
	addMetadata(xImag, "comment", "imag")
	dbg.Println("   xImag:", xImag.LLString())
	yImag := fn.cur.NewExtractValue(y, 1)
	addMetadata(yImag, "comment", "imag")
	dbg.Println("   yImag:", yImag.LLString())
	imag := op(xImag, yImag)
	addMetadata(imag, "comment", "imag")
	dbg.Println("   imag:", imag.LLString())
	// result.
	alloca := fn.cur.NewAlloca(typ)
	result := fn.cur.NewLoad(typ, alloca)
	tmp1 := fn.cur.NewInsertValue(result, real, 0)
	dbg.Println("   tmp1:", tmp1.LLString())
	tmp2 := fn.cur.NewInsertValue(result, imag, 1)
	dbg.Println("   tmp2:", tmp2.LLString())
	inst := tmp2
	return inst
}

// --- [ call instruction ] ----------------------------------------------------

// emitCall compiles the given Go SSA call instruction to corresponding LLVM IR
// instructions, emitting to fn.
func (fn *Func) emitCall(goInst *ssa.Call) error {
	dbg.Println("emitCall")
	// Function arguments; convert function arguments early, as we may rely on
	// their type when synthesizing builtin functions.
	var args []irvalue.Value
	for _, goArg := range goInst.Call.Args {
		arg := fn.useValue(goArg)
		args = append(args, arg)
	}
	// Receiver (invoke mode) or func value (call mode).
	var callee irvalue.Value
	if goCallee, ok := goInst.Call.Value.(*ssa.Builtin); ok {
		// Synthesize generic builtin `len` function based on argument type.
		switch goCallee.Name() {
		case "len":
			callee = fn.m.synthLen(args[0].Type())
		//case "cap":
		case "println":
			// TODO: synthesize `println` based on argument types?
			// implemented in builtin.ll
			callee = fn.useValue(goInst.Call.Value)
		default:
			panic(fmt.Errorf("support for builtin function %q not yet implemented", goCallee.Name()))
		}
	} else {
		callee = fn.useValue(goInst.Call.Value)
	}
	if goInst.Call.Method != nil {
		// Receiver mode.
		panic("support for receiver mode (method invocation) of Go SSA call instruction not yet implemented")
	}
	dbg.Println("   callee:", callee.Ident())
	// Bitcast pointer types of "ssa:wrapnilchk" call as follows.
	//
	//    * first argument: from T* to i8*
	//    * return value: from i8* to T*
	isWrapNilChk := false
	if named, ok := callee.(irvalue.Named); ok {
		isWrapNilChk = named.Name() == "ssa:wrapnilchk"
	}
	var tType irtypes.Type
	if isWrapNilChk {
		tType = args[0].Type()
		bitCastInst := fn.cur.NewBitCast(args[0], irtypes.I8Ptr)
		dbg.Println("   bitCastInst:", bitCastInst)
		args[0] = bitCastInst
	}
	var inst irValueInstruction = fn.cur.NewCall(callee, args...)
	if isWrapNilChk {
		bitCastInst := fn.cur.NewBitCast(inst, tType)
		dbg.Println("   bitCastInst:", bitCastInst)
		inst = bitCastInst
	}
	if !irtypes.Equal(inst.Type(), irtypes.Void) {
		inst.SetName(goInst.Name())
		fn.locals[goInst] = inst
	}
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ convert instruction ] -------------------------------------------------

// emitConvert compiles the given Go SSA convert instruction to corresponding
// LLVM IR instructions, emitting to fn.
//
// Conversions are permitted:
//    - between real numeric types.
//    - between complex numeric types.
//    - between string and []byte or []rune.
//    - between pointers and unsafe.Pointer.
//    - between unsafe.Pointer and uintptr.
//    - from (Unicode) integer to (UTF-8) string.
func (fn *Func) emitConvert(goInst *ssa.Convert) error {
	dbg.Println("emitConvert")
	from := fn.useValue(goInst.X)
	to := fn.m.irTypeFromGo(goInst.Type())
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
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ extract instruction ] -------------------------------------------------

// emitExtract compiles the given Go SSA extract instruction to corresponding
// LLVM IR instructions, emitting to fn.
func (fn *Func) emitExtract(goInst *ssa.Extract) error {
	dbg.Println("emitExtract")
	tuple := fn.useValue(goInst.Tuple)
	inst := fn.cur.NewExtractValue(tuple, uint64(goInst.Index))
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ field address instruction ] -------------------------------------------

// emitFieldAddr compiles the given Go SSA fieldaddr instruction to
// corresponding LLVM IR instructions, emitting to fn.
func (fn *Func) emitFieldAddr(goInst *ssa.FieldAddr) error {
	dbg.Println("emitFieldAddr")
	x := fn.useValue(goInst.X)
	var inst irValueInstruction
	switch xType := x.Type().(type) {
	// *struct
	case *irtypes.PointerType:
		zero := irconstant.NewInt(irtypes.I64, 0)
		// Must use i32 instead of i64 when indexing into struct types.
		//
		// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
		//
		// > The type of each index argument depends on the type it is indexing
		// > into. When indexing into a (optionally packed) structure, only i32
		// > integer constants are allowed
		field := irconstant.NewInt(irtypes.I32, int64(goInst.Field))
		indices := []irvalue.Value{
			zero,
			field,
		}
		inst = fn.cur.NewGetElementPtr(xType.ElemType, x, indices...)
	default:
		panic(fmt.Errorf("support for %T of fieldaddr instruction not yet implemented", xType))
	}
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ index address instruction ] -------------------------------------------

// emitIndexAddr compiles the given Go SSA indexaddr instruction to
// corresponding LLVM IR instructions, emitting to fn.
func (fn *Func) emitIndexAddr(goInst *ssa.IndexAddr) error {
	dbg.Println("emitIndexAddr")
	x := fn.useValue(goInst.X)
	index := fn.useValue(goInst.Index)
	var inst irValueInstruction
	switch xType := x.Type().(type) {
	// *array
	case *irtypes.PointerType:
		// TODO: Verify that index into arrays are working as intended.
		zero := irconstant.NewInt(irtypes.I64, 0)
		indices := []irvalue.Value{
			zero,
			index,
		}
		inst = fn.cur.NewGetElementPtr(xType.ElemType, x, indices...)
	case *irtypes.StructType:
		switch {
		// slice
		case strings.HasPrefix(xType.Name(), "[]"): // TODO: use other approach than type name to identify slice types; this approach is fragile as it doesn't handle type definitions (e.g. `type T []int`).
			dataType := xType.Fields[0].(*irtypes.PointerType)
			data := fn.cur.NewExtractValue(x, 0)
			addMetadata(data, "field", "data")
			dbg.Println("   data:", data.LLString())
			length := fn.cur.NewExtractValue(x, 1)
			addMetadata(length, "field", "len")
			capacity := fn.cur.NewExtractValue(x, 2)
			addMetadata(capacity, "field", "cap")
			// TODO: add bounds check of index against 0, len and cap.
			//
			//    0 <= index <= len <= cap
			// TODO: Verify that index into slices are working as intended.
			inst = fn.cur.NewGetElementPtr(dataType.ElemType, data, index)
		default:
			panic(fmt.Errorf("support for type %T (%q) in indexaddr instruction not yet implemented", xType, xType.Name()))
		}
	default:
		panic(fmt.Errorf("support for %T of indexaddr instruction not yet implemented", xType))
	}
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ lookup instruction ] --------------------------------------------------

// emitLookup compiles the given Go SSA lookup instruction to corresponding LLVM
// IR instructions, emitting to fn.
func (fn *Func) emitLookup(goInst *ssa.Lookup) error {
	dbg.Println("emitLookup")
	x := fn.useValue(goInst.X)
	index := fn.useValue(goInst.Index)
	var inst irValueInstruction
	switch xType := x.Type().(type) {
	case *irtypes.StructType:
		switch {
		// string
		case xType.Name() == "string":
			dataType := xType.Fields[0].(*irtypes.PointerType)
			data := fn.cur.NewExtractValue(x, 0)
			addMetadata(data, "field", "data")
			gep := fn.cur.NewGetElementPtr(dataType.ElemType, data, index)
			dbg.Println("   gep:", gep.LLString())
			inst = fn.cur.NewLoad(dataType.ElemType, gep)
		// map
		case strings.HasPrefix(xType.Name(), "map["): // TODO: use other approach than type name to identify map types; this approach is fragile as it doesn't handle type definitions (e.g. `type T map[int]string`).
			if goInst.CommaOk {
				// TODO: add support for `comma, ok` use of map.
			}
			panic(fmt.Errorf("support for type %T (%q) in lookup instruction not yet implemented", xType, xType.Name()))
		default:
			panic(fmt.Errorf("support for type %T (%q) in lookup instruction not yet implemented", xType, xType.Name()))
		}
	default:
		panic(fmt.Errorf("support for type %T (%q) in lookup instruction not yet implemented", xType, xType.Name()))
	}
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ phi instruction ] -----------------------------------------------------

// emitPhi compiles the given Go SSA phi instruction to corresponding LLVM IR
// instructions, emitting to fn.
func (fn *Func) emitPhi(goInst *ssa.Phi) error {
	dbg.Println("emitPhi")
	var incs []*ir.Incoming
	for i, goEdge := range goInst.Edges {
		x := fn.useValue(goEdge)
		goPred := goInst.Block().Preds[i]
		pred := fn.getBlock(goPred)
		inc := ir.NewIncoming(x, pred)
		incs = append(incs, inc)
	}
	inst := fn.cur.NewPhi(incs...)
	if len(goInst.Comment) > 0 {
		addMetadata(inst, "comment", goInst.Comment)
	}
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ slice instruction ] ---------------------------------------------------

// emitSlice compiles the given Go SSA slice instruction to corresponding LLVM
// IR instructions, emitting to fn.
func (fn *Func) emitSlice(goInst *ssa.Slice) error {
	dbg.Println("emitSlice")
	x := fn.useValue(goInst.X) // slice, string, or *array
	dbg.Println("   x:", x.String())
	var low, high, max irvalue.Value
	if goInst.Low != nil {
		low = fn.useValue(goInst.Low)
	}
	if goInst.High != nil {
		high = fn.useValue(goInst.High)
	}
	if goInst.Max != nil {
		max = fn.useValue(goInst.Max)
	}
	// Get element type.
	var elemType irtypes.Type
	var (
		data     irvalue.Value
		length   irvalue.Value
		capacity irvalue.Value
	)
	switch xType := x.Type().(type) {
	// slice, string
	case *irtypes.StructType:
		switch {
		// slice
		case strings.HasPrefix(xType.Name(), "[]"): // TODO: use other approach than type name to identify slice types; this approach is fragile as it doesn't handle type definitions (e.g. `type T []int`).
			dataType := xType.Fields[0].(*irtypes.PointerType)
			elemType = dataType.ElemType
			dataField := fn.cur.NewExtractValue(x, 0)
			addMetadata(dataField, "field", "data")
			data = dataField
			lengthField := fn.cur.NewExtractValue(x, 1)
			addMetadata(lengthField, "field", "len")
			length = lengthField
			capacityField := fn.cur.NewExtractValue(x, 2)
			addMetadata(capacityField, "field", "cap")
			capacity = capacityField
		// string
		case xType.Name() == "string":
			elemType = fn.m.irTypeFromName("uint8") // TODO: use byte alias instead of uint8.
			dataField := fn.cur.NewExtractValue(x, 0)
			addMetadata(dataField, "field", "data")
			data = dataField
			lengthField := fn.cur.NewExtractValue(x, 1)
			addMetadata(lengthField, "field", "len")
			length = lengthField
			capacity = lengthField
		default:
			panic(fmt.Errorf("support for type %T (%q) in slice instruction not yet implemented", xType, xType.Name()))
		}
	// *array
	case *irtypes.PointerType:
		arrayType, ok := xType.ElemType.(*irtypes.ArrayType)
		if !ok {
			panic(fmt.Errorf("support for type %T (%v) in slice instruction not yet implemented", xType, xType.String()))
		}
		elemType = arrayType.ElemType
		zero := irconstant.NewInt(irtypes.I64, 0)
		indices := []irvalue.Value{
			zero,
			zero,
		}
		dataField := fn.cur.NewGetElementPtr(arrayType, x, indices...)
		addMetadata(dataField, "field", "data")
		data = dataField
		length = irconstant.NewInt(fn.m.irTypeFromName("int").(*irtypes.IntType), int64(arrayType.Len))
		capacity = length
	default:
		panic(fmt.Errorf("support for type %T in slice instruction not yet implemented", xType))
	}
	// Allocate new slice value.
	sliceType := fn.m.newSliceType(elemType)
	alloca := fn.cur.NewAlloca(sliceType)
	fn.cur.NewStore(irconstant.NewZeroInitializer(sliceType), alloca)
	var slice irvalue.Value = fn.cur.NewLoad(sliceType, alloca)
	// TODO: add bounds check of low, high and max.
	// data[low::]
	if low != nil {
		data = fn.cur.NewGetElementPtr(elemType, data, low)
	}
	insertData := fn.cur.NewInsertValue(slice, data, 0)
	addMetadata(insertData, "field", "data")
	slice = insertData
	// data[:high:]
	if high != nil {
		length = high
	}
	insertLength := fn.cur.NewInsertValue(slice, length, 1)
	addMetadata(insertLength, "field", "len")
	slice = insertLength
	// data[::max]
	if max != nil {
		capacity = max
	}
	insertCapacity := fn.cur.NewInsertValue(slice, capacity, 2)
	addMetadata(insertCapacity, "field", "cap")
	inst := insertCapacity
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}

// --- [ unary operation instruction ] -----------------------------------------

// emitUnOp compiles the given Go SSA unary operation instruction to
// corresponding LLVM IR instructions, emitting to fn.
func (fn *Func) emitUnOp(goInst *ssa.UnOp) error {
	dbg.Println("emitUnOp")
	dbg.Println("   op:", goInst.Op)
	x := fn.useValue(goInst.X)
	dbg.Println("   x:", x)
	var inst irValueInstruction
	switch goInst.Op {
	// Logical negation.
	case token.NOT: // !
		inst = fn.cur.NewXor(x, irconstant.True)
	// Negation.
	case token.SUB: // -
		// Note that the `sub` instruction is used to represent the `neg`
		// instruction present in most other intermediate representations.
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			zero := irconstant.NewInt(typ, 0)
			inst = fn.cur.NewSub(zero, x)
		case *irtypes.FloatType:
			zero := irconstant.NewFloat(typ, 0)
			inst = fn.cur.NewFSub(zero, x)
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	// Channel receive.
	case token.ARROW: // <-
		if goInst.CommaOk {
			// The result is a 2-tuple of the value and a boolean indicating the
			// success of the receive. The components of the tuple are accessed
			// using Extract.
			panic("support for comma-ok of Go SSA unary operation instruction token ARROW (<-) not yet implemented")
		}
		panic("support for Go SSA unary operation instruction token ARROW (<-) not yet implemented")
	// Pointer indirection (load).
	case token.MUL: // *
		elemType := x.Type().(*irtypes.PointerType).ElemType
		inst = fn.cur.NewLoad(elemType, x)
	// Bitwise complement.
	case token.XOR: // ^
		switch typ := x.Type().(type) {
		case *irtypes.IntType:
			zero := irconstant.NewInt(typ, 0)
			inst = fn.cur.NewXor(x, zero)
		default:
			panic(fmt.Errorf("support for operand type %T (%q) of Go SSA binary operation instruction (%v) not yet implemented", typ, typ.Name(), goInst.Op))
		}
	default:
		panic(fmt.Errorf("support for Go SSA unary operation instruction token %v not yet implemented", goInst.Op))
	}
	inst.SetName(goInst.Name())
	// Add unary operation token metadata attachment to instruction.
	addMetadata(inst, "unary_op", goInst.Op.String())
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
}
