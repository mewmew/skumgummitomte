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

	// Entry basic block of function, used for allocation of local variables.
	Entry *ir.Block
	// Current basic block being generated.
	Cur *ir.Block
	// Maps from LLVM IR basic block name to LLVM IR basic block, with one LLVM
	// IR basic block per Go SSA basic block.
	blocks map[string]*ir.Block
	// Locals maps from local Go SSA variable name to LLVM IR value of the
	// corresponding LLVM IR local variable.
	locals map[string]*ir.InstAlloca
}

// NewFunc returns a new LLVM IR function generator based on the given LLVM IR
// function declaration.
func NewFunc(f *ir.Func) *Func {
	entry := f.NewBlock("entry")
	return &Func{
		Func:   f,
		Entry:  entry,
		Cur:    entry,
		blocks: make(map[string]*ir.Block),
		locals: make(map[string]*ir.InstAlloca),
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

	// Locate output LLVM IR function.
	f, ok := gen.funcs[funcName]
	if !ok {
		return errors.Errorf("unable to locate LLVM IR function %q", funcName)
	}

	// Index Go SSA basic blocks by creating corresponding LLVM IR basic blocks.
	fn := NewFunc(f)
	for _, goBlock := range goFunc.Blocks {
		blockName := getBlockName(goBlock.Index)
		block := ir.NewBlock(blockName)
		fn.blocks[blockName] = block
	}

	// Add unconditional branch from LLVM IR entry basic block to Go SSA entry
	// basic block.
	entryBlock := fn.getBlock(0)
	fn.Cur.NewBr(entryBlock)

	// Generate LLVM IR basic blocks of Go function definition.
	//
	// Process basic blocks in dominance order, starting with dominators and
	// sorting equal dominance by basic block index.
	// TODO: sort blocks by dom.
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
	block := fn.getBlock(goBlock.Index)
	fn.Func.Blocks = append(fn.Func.Blocks, block)
	fn.Cur = block
	for _, goInst := range goBlock.Instrs {
		if err := gen.compileInst(goInst, fn); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
