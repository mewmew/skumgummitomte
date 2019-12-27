package main

import (
	"fmt"

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
	case *ssa.Return:
		return gen.compileReturn(goInst, fn)
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
	_ = callInst
	return nil
}

// --- [ return instruction ] --------------------------------------------------

// compileReturn compiles the given Go SSA return instruction into corresponding
// LLVM IR instructions, emitting to fn.
func (gen *generator) compileReturn(goInst *ssa.Return, fn *Func) error {
	// TODO: support return values.
	fn.Cur.NewRet(nil)
	return nil
}
