package irgen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// Func is an LLVM IR function generator.
type Func struct {
	// Output LLVM IR function.
	*ir.Func
	// Input Go SSA function.
	goFunc *ssa.Function

	// Parent LLVM IR module generator.
	m *Module
	// Entry basic block of function, used for allocation of local variables.
	entry *ir.Block
	// Current basic block being generated.
	cur *ir.Block
	// Maps from function local Go SSA value to corresponding LLVM IR value in
	// the LLVM IR function being generated.
	locals map[ssa.Value]irvalue.Value
	// Maps from Go SSA basic block to corresponding LLVM IR basic block in the
	// LLVM IR function being generated.
	blocks map[*ssa.BasicBlock]*ir.Block
}

// NewFunc returns a new LLVM IR function generator for the given Go SSA
// function, emitting to m.
//
// Pre-condition: index functions of m.
func (m *Module) NewFunc(goFunc *ssa.Function) *Func {
	f := m.getFunc(goFunc)
	entry := f.NewBlock("entry")
	return &Func{
		Func:   f,
		goFunc: goFunc,
		m:      m,
		entry:  entry,
		cur:    entry,
		locals: make(map[ssa.Value]irvalue.Value),
		blocks: make(map[*ssa.BasicBlock]*ir.Block),
	}
}

// --- [ init ] ----------------------------------------------------------------

// initPredeclaredFuncs initializes LLVM IR functions corresponding to the
// predeclared functions in Go (e.g. "println").
//
// Pre-condition: initialize predeclared types.
func (m *Module) initPredeclaredFuncs() {
	// println.
	retType := irtypes.Void
	param := ir.NewParam("", m.irTypeFromName("string"))
	printlnFunc := m.Module.NewFunc("println", retType, param)
	printlnFunc.Sig.Variadic = true
	m.predeclaredFuncs[printlnFunc.Name()] = printlnFunc
}

// --- [ get ] -----------------------------------------------------------------

// getFunc returns the LLVM IR function corresponding to the given Go SSA
// function.
//
// Pre-condition: index functions of m.
func (m *Module) getFunc(goFunc *ssa.Function) *ir.Func {
	// Lookup indexed LLVM IR function declaration of Go SSA function.
	global, ok := m.globals[goFunc]
	if !ok {
		// Pre-condition invalidated, function declaration not indexed. This is a
		// fatal error and indicates a bug in irgen.
		panic(fmt.Errorf("unable to locate indexed LLVM IR function declaration of Go SSA function %q", goFunc.Name()))
	}
	return global.(*ir.Func)
}

// --- [ index ] ---------------------------------------------------------------

// indexFunc indexes the given Go SSA function, creating a corresponding LLVM IR
// function, emitting to m.
func (m *Module) indexFunc(goFunc *ssa.Function) error {
	// Convert Go function parameters to equivalent LLVM IR function parameters.
	var params []*ir.Param
	goParams := goFunc.Signature.Params()
	for i := 0; i < goParams.Len(); i++ {
		goParam := goParams.At(i)
		paramName := goParam.Name()
		paramType := m.irTypeFromGo(goParam.Type())
		param := ir.NewParam(paramName, paramType)
		params = append(params, param)
	}
	// Convert Go function return types to equivalent LLVM IR function return
	// types.
	var resultTypes []irtypes.Type
	goResults := goFunc.Signature.Results()
	for i := 0; i < goResults.Len(); i++ {
		goResult := goResults.At(i)
		resultName := goResult.Name()
		// TODO: add resultName as field name of (custom) result structure type.
		_ = resultName
		resultType := m.irTypeFromGo(goResult.Type())
		resultTypes = append(resultTypes, resultType)
	}
	// Convert multiple return types a single return type by creating a structure
	// type with one field per return type.
	var retType irtypes.Type
	switch len(resultTypes) {
	// void return.
	case 0:
		retType = irtypes.Void
	// single return type.
	case 1:
		retType = resultTypes[0]
	// multiple return types.
	default:
		retType = irtypes.NewStruct(resultTypes...)
	}
	// Generate LLVM IR function declaration, emitting to m.
	f := m.Module.NewFunc(goFunc.Name(), retType, params...)
	f.Sig.Variadic = goFunc.Signature.Variadic()
	// Index LLVM IR function declaration.
	m.globals[goFunc] = f
	return nil
}

// --- [ compile ] -------------------------------------------------------------

// emitFunc compiles the given Go SSA function into LLVM IR, emitting to m.
//
// Pre-condition: index functions of m.
func (m *Module) emitFunc(goFunc *ssa.Function) error {
	// Early return for function declaration (without body).
	if len(goFunc.Blocks) == 0 {
		return nil
	}
	// Index Go SSA basic blocks by creating corresponding LLVM IR basic blocks.
	fn := m.NewFunc(goFunc)
	for _, goBlock := range goFunc.Blocks {
		if err := fn.indexBlock(goBlock); err != nil {
			return errors.WithStack(err)
		}
	}
	// Add unconditional branch from LLVM IR entry basic block to Go SSA entry
	// basic block.
	entryBlock := fn.getBlock(goFunc.Blocks[0])
	fn.cur.NewBr(entryBlock)
	// Generate LLVM IR basic blocks of Go function definition.
	//
	// Process basic blocks in dominance order, starting with dominators and
	// sorting equal dominance by basic block index.
	// TODO: sort blocks by dom.
	for _, goBlock := range goFunc.Blocks {
		if err := fn.emitBlock(goBlock); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
