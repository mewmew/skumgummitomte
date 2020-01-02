package irgen

import (
	"fmt"
	"go/token"

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
		panic("support for *ssa.DebugRef not yet implemented")
	case *ssa.Defer:
		panic("support for *ssa.Defer not yet implemented")
	case *ssa.Go:
		panic("support for *ssa.Go not yet implemented")
	case *ssa.If:
		return fn.emitIf(goInst)
	case *ssa.Jump:
		return fn.emitJump(goInst)
	case *ssa.MapUpdate:
		panic("support for *ssa.MapUpdate not yet implemented")
	case *ssa.Panic:
		panic("support for *ssa.Panic not yet implemented")
	case *ssa.Return:
		return fn.emitReturn(goInst)
	case *ssa.RunDefers:
		// TODO: implement support for defer.
		return nil // ignore *ssa.RunDefers instruction for now
	case *ssa.Send:
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
		panic("support for *ssa.ChangeInterface not yet implemented")
	case *ssa.ChangeType:
		panic("support for *ssa.ChangeType not yet implemented")
	case *ssa.Convert:
		panic("support for *ssa.Convert not yet implemented")
	case *ssa.Extract:
		panic("support for *ssa.Extract not yet implemented")
	case *ssa.Field:
		panic("support for *ssa.Field not yet implemented")
	case *ssa.FieldAddr:
		panic("support for *ssa.FieldAddr not yet implemented")
	case *ssa.Index:
		panic("support for *ssa.Index not yet implemented")
	case *ssa.IndexAddr:
		panic("support for *ssa.IndexAddr not yet implemented")
	case *ssa.Lookup:
		panic("support for *ssa.Lookup not yet implemented")
	case *ssa.MakeChan:
		panic("support for *ssa.MakeChan not yet implemented")
	case *ssa.MakeClosure:
		panic("support for *ssa.MakeClosure not yet implemented")
	case *ssa.MakeInterface:
		panic("support for *ssa.MakeInterface not yet implemented")
	case *ssa.MakeMap:
		panic("support for *ssa.MakeMap not yet implemented")
	case *ssa.MakeSlice:
		panic("support for *ssa.MakeSlice not yet implemented")
	case *ssa.Next:
		panic("support for *ssa.Next not yet implemented")
	case *ssa.Phi:
		return fn.emitPhi(goInst)
	case *ssa.Range:
		panic("support for *ssa.Range not yet implemented")
	case *ssa.Select:
		panic("support for *ssa.Select not yet implemented")
	case *ssa.Slice:
		panic("support for *ssa.Slice not yet implemented")
	case *ssa.TypeAssert:
		panic("support for *ssa.TypeAssert not yet implemented")
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
		ret := fn.entry.NewAlloca(structType)
		for index, result := range results {
			fn.cur.NewInsertValue(ret, result, uint64(index))
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
	dbg.Println("   inst:", inst)
	// The space of a stack allocated local variable is re-initialized to zero
	// each time it is executed.
	zero := irconstant.NewZeroInitializer(ptrType.ElemType)
	fn.cur.NewStore(zero, inst)
	return nil
}

// ~~~ [ new - heap alloc instruction ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// emitNew compiles the given Go SSA heap alloc instruction to corresponding
// LLVM IR instructions, emitting to fn.
func (fn *Func) emitNew(goInst *ssa.Alloc) error {
	dbg.Println("emitNew")
	typ := fn.m.irTypeFromGo(goInst.Type())
	ptrType := typ.(*irtypes.PointerType)
	// Define `new(T)` function if not present.
	typeName := ptrType.ElemType.Name() // TODO: check that builtin type names (e.g. `int32`) are currently handled.
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
		addMetadata(inst, "var_name", goInst.Comment)
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
		}
	// EQL (==)
	case token.EQL: // ==
		panic("support for Go SSA binary operation instruction token EQL (==) not yet implemented")
	// NEQ (!=)
	case token.NEQ: // !=
		panic("support for Go SSA binary operation instruction token NEQ (!=) not yet implemented")
	// LSS (<)
	case token.LSS: // <
		panic("support for Go SSA binary operation instruction token LSS (<) not yet implemented")
	// LEQ (<=)
	case token.LEQ: // <=
		panic("support for Go SSA binary operation instruction token LEQ (<=) not yet implemented")
	// GTR (<)
	case token.GTR: // <
		panic("support for Go SSA binary operation instruction token GTR (<) not yet implemented")
	// GEQ (>=)
	case token.GEQ: // >=
		panic("support for Go SSA binary operation instruction token GEQ (>=) not yet implemented")
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
	// Receiver (invoke mode) or func value (call mode).
	callee := fn.useValue(goInst.Call.Value)
	if goInst.Call.Method != nil {
		// Receiver mode.
		panic("support for receiver mode (method invocation) of Go SSA call instruction not yet implemented")
	}
	dbg.Println("   callee:", callee.Ident())
	// Function arguments.
	var args []irvalue.Value
	for _, goArg := range goInst.Call.Args {
		arg := fn.useValue(goArg)
		args = append(args, arg)
	}
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

// --- [ phi instruction ] -----------------------------------------------------

// emitPhi compiles the given Go SSA phi instruction to corresponding LLVM IR
// instructions, emitting to fn.
func (fn *Func) emitPhi(goInst *ssa.Phi) error {
	dbg.Println("emitPhi")
	var incs []*ir.Incoming
	for i, goEdge := range goInst.Edges {
		x := fn.irValueFromGo(goEdge)
		goPred := goInst.Block().Preds[i]
		pred := fn.getBlock(goPred)
		inc := ir.NewIncoming(x, pred)
		incs = append(incs, inc)
	}
	inst := fn.cur.NewPhi(incs...)
	if len(goInst.Comment) > 0 {
		addMetadata(inst, "comment", goInst.Comment)
	}
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
			panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", typ, goInst.Op))
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
