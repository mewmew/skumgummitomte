package irgen

import (
	"fmt"
	"go/token"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/metadata"
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
		panic("support for *ssa.Phi not yet implemented")
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
		fmt.Println("   result:", result)
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
		//  Allocate space in the heap.
		panic("support for heap allocated space of Go SSA alloc instruction not yet implemented")
	}
	typ := fn.m.irTypeFromGo(goInst.Type())
	ptrType := typ.(*irtypes.PointerType)
	inst := fn.cur.NewAlloca(ptrType.ElemType)
	inst.SetName(goInst.Name())
	fn.locals[goInst] = inst
	// Add local variable name metadata attachment to alloca instruction.
	if len(goInst.Comment) > 0 {
		mdLocalName := &metadata.Attachment{
			Name: "var_name",
			Node: &metadata.Tuple{
				MetadataID: -1, // metadata literal.
				Fields: []metadata.Field{
					&metadata.String{Value: goInst.Comment},
				},
			},
		}
		inst.Metadata = append(inst.Metadata, mdLocalName)
	}
	dbg.Println("   inst:", inst)
	// The space of a stack allocated local variable is re-initialized to zero
	// each time it is executed.
	zero := irconstant.NewZeroInitializer(ptrType.ElemType)
	fn.cur.NewStore(zero, inst)
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
	var isInt, isFloat bool
	switch x.Type().(type) {
	case *irtypes.IntType:
		isInt = true
	case *irtypes.FloatType:
		isFloat = true
	default:
		panic(fmt.Errorf("support for operand type %T of Go SSA binary operation instruction (%v) not yet implemented", x.Type(), goInst.Op))
	}
	var inst irValueInstruction
	switch goInst.Op {
	// ADD (+)
	case token.ADD: // +
		switch {
		case isInt:
			inst = fn.cur.NewAdd(x, y)
		case isFloat:
			inst = fn.cur.NewFAdd(x, y)
		}
	// SUB (-)
	case token.SUB: // -
		panic("support for Go SSA binary operation instruction token SUB (-) not yet implemented")
	// MUL (*)
	case token.MUL: // *
		panic("support for Go SSA binary operation instruction token MUL (*) not yet implemented")
	// QUO (/)
	case token.QUO: // /
		panic("support for Go SSA binary operation instruction token QUO (/) not yet implemented")
	// REM (%)
	case token.REM: // %
		panic("support for Go SSA binary operation instruction token REM (%) not yet implemented")
	// AND (&)
	case token.AND: // &
		panic("support for Go SSA binary operation instruction token AND (&) not yet implemented")
	// OR (|)
	case token.OR: // |
		panic("support for Go SSA binary operation instruction token OR (|) not yet implemented")
	// XOR (^)
	case token.XOR: // ^
		panic("support for Go SSA binary operation instruction token XOR (^) not yet implemented")
	// SHL (<<)
	case token.SHL: // <<
		panic("support for Go SSA binary operation instruction token SHL (<<) not yet implemented")
	// SHR (>>)
	case token.SHR: // >>
		panic("support for Go SSA binary operation instruction token SHR (>>) not yet implemented")
	// AND_NOT (&^)
	case token.AND_NOT: // &^
		panic("support for Go SSA binary operation instruction token AND_NOT (&^) not yet implemented")
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
	fn.locals[goInst] = inst
	dbg.Println("   inst:", inst.LLString())
	return nil
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

// --- [ unary operation instruction ] -----------------------------------------

// emitUnOp compiles the given Go SSA unary operation instruction to
// corresponding LLVM IR instructions, emitting to fn.
func (fn *Func) emitUnOp(goInst *ssa.UnOp) error {
	dbg.Println("emitUnOp")
	dbg.Println("   op:", goInst.Op)
	x := fn.useValue(goInst.X)
	dbg.Println("   x:", x)
	switch goInst.Op {
	// Logical negation.
	case token.NOT: // !
		panic("support for Go SSA unary operation instruction token NOT (!) not yet implemented")
	// Negation.
	case token.SUB: // -
		panic("support for Go SSA unary operation instruction token SUB (-) not yet implemented")
	// Channel receive.
	case token.ARROW: // <-
		if goInst.CommaOk {
			// The result is a 2-tuple of the value and a boolean indicating the
			// success of the receive. The components of the tuple are accessed
			// using Extract.
		}
		panic("support for Go SSA unary operation instruction token ARROW (<-) not yet implemented")
	// Pointer indirection (load).
	case token.MUL: // *
		elemType := x.Type().(*irtypes.PointerType).ElemType
		inst := fn.cur.NewLoad(elemType, x)
		inst.SetName(goInst.Name())
		fn.locals[goInst] = inst
		dbg.Println("   inst:", inst.LLString())
		return nil
	// Bitwise complement.
	case token.XOR: // ^
		panic("support for Go SSA unary operation instruction token XOR (^) not yet implemented")
	default:
		panic(fmt.Errorf("support for Go SSA unary operation instruction token %v not yet implemented", goInst.Op))
	}
}
