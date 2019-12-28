package main

import (
	"fmt"
	"go/token"

	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// compileInst compiles the given Go SSA instruction into corresponding LLVM IR
// instructions, emitting to fn.
func (gen *generator) compileInst(goInst ssa.Instruction, fn *Func) error {
	switch goInst := goInst.(type) {
	case ssa.CallInstruction:
		return gen.compileCallInst(goInst, fn)
	case *ssa.If:
		return gen.compileIf(goInst, fn)
	case *ssa.Jump:
		return gen.compileJump(goInst, fn)
	case *ssa.Return:
		return gen.compileReturn(goInst, fn)
	case *ssa.RunDefers:
		// TODO: implement support for defer.
		return nil // ignore *ssa.RunDefers instruction for now
	case *ssa.Store:
		return gen.compileStore(goInst, fn)
	case *ssa.UnOp:
		return gen.compileUnOp(goInst, fn)
	default:
		panic(fmt.Errorf("support for Go SSA instruction %T not yet implemented", goInst))
	}
}

// --- [ call instruction ] ----------------------------------------------------

// compileCallInst compiles the given Go SSA call instruction into corresponding
// LLVM IR instructions, emitting to fn.
func (gen *generator) compileCallInst(goInst ssa.CallInstruction, fn *Func) error {
	switch goInst := goInst.(type) {
	case *ssa.Call:
		return gen.compileCall(goInst, fn)
	default:
		panic(fmt.Errorf("support for Go SSA instruction %T not yet implemented", goInst))
	}
}

// compileCall compiles the given Go SSA call instruction into corresponding
// LLVM IR instructions, emitting to fn.
func (gen *generator) compileCall(goInst *ssa.Call, fn *Func) error {
	// Convert Go callee to an equivalent LLVM IR callee.
	var callee value.Value
	switch goCallee := goInst.Call.Value.(type) {
	case *ssa.Builtin:
		calleeFunc, ok := gen.predeclaredFuncs[goCallee.Name()]
		if !ok {
			return errors.Errorf("unable to locate LLVM IR function of predeclared Go function %q", goCallee.Name())
		}
		dbg.Println("callee:", calleeFunc.LLString())
		callee = calleeFunc
	default:
		panic(fmt.Errorf("support for Go SSA call instruction callee %T not yet implemented", goCallee))
	}
	// Convert Go call arguments to an equivalent LLVM IR call arguments.
	var args []value.Value
	for _, goArg := range goInst.Call.Args {
		arg := gen.llValueFromGoValue(goArg)
		args = append(args, arg)
	}
	callInst := fn.Cur.NewCall(callee, args...)
	dbg.Println("call inst", callInst.LLString())
	if !irtypes.Equal(callInst.Type(), irtypes.Void) {
		fn.defLocal(goInst.Name(), callInst)
	}
	return nil
}

// --- [ if instruction ] ------------------------------------------------------

// compileIf compiles the given Go SSA if instruction into corresponding LLVM IR
// instructions, emitting to fn.
func (gen *generator) compileIf(goInst *ssa.If, fn *Func) error {
	// The If instruction transfers control to one of the two successors of its
	// owning block, depending on the boolean Cond: the first if true, the second
	// if false.
	cond := gen.useValue(goInst.Cond, fn)
	succs := goInst.Block().Succs
	targetTrue := fn.getBlock(succs[0].Index)
	targetFalse := fn.getBlock(succs[1].Index)
	condBrTerm := fn.Cur.NewCondBr(cond, targetTrue, targetFalse)
	dbg.Println("cond br term:", condBrTerm.LLString())
	return nil
}

// --- [ jump instruction ] ----------------------------------------------------

// compileJump compiles the given Go SSA jump instruction into corresponding
// LLVM IR instructions, emitting to fn.
func (gen *generator) compileJump(goInst *ssa.Jump, fn *Func) error {
	// The Jump instruction transfers control to the sole successor of its owning
	// block.
	succs := goInst.Block().Succs
	target := fn.getBlock(succs[0].Index)
	brTerm := fn.Cur.NewBr(target)
	dbg.Println("br term:", brTerm.LLString())
	return nil
}

// --- [ return instruction ] --------------------------------------------------

// compileReturn compiles the given Go SSA return instruction into corresponding
// LLVM IR instructions, emitting to fn.
func (gen *generator) compileReturn(goInst *ssa.Return, fn *Func) error {
	// TODO: support return values.
	retTerm := fn.Cur.NewRet(nil)
	dbg.Println("ret term:", retTerm.LLString())
	return nil
}

// --- [ store instruction ] ---------------------------------------------------

// compileStore compiles the given Go SSA store instruction into corresponding
// LLVM IR instructions, emitting to fn.
func (gen *generator) compileStore(goInst *ssa.Store, fn *Func) error {
	src := gen.useValue(goInst.Val, fn)
	dst := gen.useValue(goInst.Addr, fn)
	storeInst := fn.Cur.NewStore(src, dst)
	dbg.Println("store inst:", storeInst.LLString())
	return nil
}

// --- [ unary operation instruction ] -----------------------------------------

// compileUnOp compiles the given Go SSA unary operation instruction into
// corresponding LLVM IR instructions, emitting to fn.
func (gen *generator) compileUnOp(goInst *ssa.UnOp, fn *Func) error {
	switch goInst.Op {
	// Logical negation.
	case token.NOT:
		panic(fmt.Errorf("support for Go SSA unary operation instruction token NOT (!) not yet implemented"))
	// Negation.
	case token.SUB:
		panic(fmt.Errorf("support for Go SSA unary operation instruction token SUB (-) not yet implemented"))
	// Channel receive.
	case token.ARROW:
		if goInst.CommaOk {
			// The result is a 2-tuple of the value and a boolean indicating the
			// success of the receive. The components of the tuple are accessed
			// using Extract.
		}
		panic(fmt.Errorf("support for Go SSA unary operation instruction token ARROW (<-) not yet implemented"))
	// Pointer indirection (load).
	case token.MUL:
		x := gen.useValue(goInst.X, fn)
		elemType := x.Type().(*irtypes.PointerType).ElemType
		loadInst := fn.Cur.NewLoad(elemType, x)
		fn.defLocal(goInst.Name(), loadInst)
		return nil
	// Bitwise complement.
	case token.XOR:
		panic(fmt.Errorf("support for Go SSA unary operation instruction token XOR (^) not yet implemented"))
	default:
		panic(fmt.Errorf("support for Go SSA unary operation instruction token %v not yet implemented", goInst.Op))
	}
}
