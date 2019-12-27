package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// initPredeclaredFuncs initializes LLVM IR functions corresponding to the
// predeclared functions in Go (e.g. "println").
//
// pre-condition: initPredeclaredTypes
func (gen *generator) initPredeclaredFuncs() {
	// println.
	// TODO: fix type of println; parameters should support arbitrary types.
	retType := irtypes.Void
	param := ir.NewParam("", gen.llTypeFromName("string"))
	printlnFunc := gen.module.NewFunc("println", retType, param)
	printlnFunc.Sig.Variadic = true
	gen.predeclaredFuncs[printlnFunc.Name()] = printlnFunc
}

// llFuncFromName returns the LLVM IR function for the corresponding Go function
// name.
func (gen *generator) llFuncFromName(typeName string) irtypes.Type {
	// TODO: handle shadowing of predeclared types based on scope.
	if typ, ok := gen.predeclaredTypes[typeName]; ok {
		return typ
	}
	if typ, ok := gen.types[typeName]; ok {
		return typ
	}
	panic(fmt.Errorf("unable to locate LLVM IR type for the corresponding Go type name %q", typeName))
}

// Func is an LLVM IR function generator.
type Func struct {
	// Output LLVM IR function.
	*ir.Func

	// Current basic block being generated.
	Cur *ir.Block
}

// NewFunc returns a new LLVM IR function generator based on the given LLVM IR
// function declaration.
func NewFunc(f *ir.Func) *Func {
	return &Func{
		Func: f,
		Cur:  f.NewBlock("entry"),
	}
}

// compileFunction compiles the given Go SSA function into LLVM IR.
func (gen *generator) compileFunction(funcName string, goFunc *ssa.Function) error {
	// TODO: remove debug output.
	dbg.Println("compileFunction")
	dbg.Println(goFunc)
	dbg.Println()

	// Early return for function declaration (without body).
	if len(goFunc.Blocks) == 0 {
		return nil
	}

	// TODO: remove once support for Go `init` functions is added.
	if funcName != "main" {
		return nil
	}

	// Locate output LLVM IR function.
	f, ok := gen.funcs[funcName]
	if !ok {
		return errors.Errorf("unable to locate LLVM IR function %q", funcName)
	}

	// Generate LLVM IR basic blocks of Go function definition.
	fn := NewFunc(f)
	for _, goBlock := range goFunc.Blocks {
		if err := gen.compileBlock(goBlock, fn); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// compileBlock compiles the given Go SSA function basic block into
// corresponding LLVM IR basic blocks, emitting to fn.
func (gen *generator) compileBlock(goBlock *ssa.BasicBlock, fn *Func) error {
	for _, goInst := range goBlock.Instrs {
		if err := gen.compileInst(goInst, fn); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
