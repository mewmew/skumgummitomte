package irgen

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
	// --- [ builtin Go functions ] ---

	// println.
	{
		retType := irtypes.Void
		param := ir.NewParam("", m.irTypeFromName("string"))
		printlnFunc := m.Module.NewFunc("println", retType, param)
		printlnFunc.Sig.Variadic = true
		m.predeclaredFuncs[printlnFunc.Name()] = printlnFunc
	}

	// --- [ needed by Go SSA code ] ---

	// wrapnilchk returns ptr if non-nil, panics otherwise.
	// (For use in indirection wrappers.)
	//
	//    func ssa:wrapnilchk(ptr *T, recvType, methodName string) *T
	{
		ptrType := irtypes.I8Ptr // generic pointer type.
		retType := ptrType
		params := []*ir.Param{
			ir.NewParam("ptr", ptrType),
			ir.NewParam("recvType", m.irTypeFromName("string")),
			ir.NewParam("methodName", m.irTypeFromName("string")),
		}
		wrapnilchkFunc := m.Module.NewFunc("ssa:wrapnilchk", retType, params...)
		m.predeclaredFuncs[wrapnilchkFunc.Name()] = wrapnilchkFunc
	}

	// --- [ dependencies of new(T) ] ---

	// calloc.
	{
		// void *calloc(size_t nmemb, size_t size)
		retType := irtypes.I8Ptr // generic pointer type.
		params := []*ir.Param{
			ir.NewParam("nmemb", m.irTypeFromName("uint64")),
			ir.NewParam("size", m.irTypeFromName("uint64")),
		}
		callocFunc := m.Module.NewFunc("calloc", retType, params...)
		m.predeclaredFuncs[callocFunc.Name()] = callocFunc
	}

	// llvm.objectsize.i64
	{
		// declare i64 @llvm.objectsize.i64(i8* <object>, i1 <min>, i1 <nullunknown>, i1 <dynamic>)
		retType := irtypes.I64
		params := []*ir.Param{
			ir.NewParam("object", irtypes.I8Ptr),
			ir.NewParam("min", irtypes.I1),
			ir.NewParam("nullunknown", irtypes.I1),
			ir.NewParam("dynamic", irtypes.I1),
		}
		objectsizeFunc := m.Module.NewFunc("llvm.objectsize.i64", retType, params...)
		m.predeclaredFuncs[objectsizeFunc.Name()] = objectsizeFunc
	}
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
		panic(fmt.Errorf("unable to locate indexed LLVM IR function declaration of Go SSA function %q", m.fullName(goFunc)))
	}
	return global.(*ir.Func)
}

// getPredeclaredFunc returns the predeclared LLVM IR function of the given
// function name.
//
// Pre-condition: initialize predeclared functions of m.
func (m *Module) getPredeclaredFunc(funcName string) *ir.Func {
	predeclaredFunc, ok := m.predeclaredFuncs[funcName]
	if !ok {
		panic(fmt.Errorf("unable to locate predeclared LLVM IR function %q", funcName))
	}
	return predeclaredFunc
}

// --- [ index ] ---------------------------------------------------------------

// indexFunc indexes the given Go SSA function, creating a corresponding LLVM IR
// function, emitting to m.
func (m *Module) indexFunc(goFunc *ssa.Function) error {
	// TODO: use m.irTypeFromGo(goFunc.Signature) to simplify m.indexFunc.
	// Convert Go function parameters to equivalent LLVM IR function parameters
	// (including receiver of methods).
	params := make([]*ir.Param, 0, len(goFunc.Params))
	for _, goParam := range goFunc.Params {
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
	f := m.Module.NewFunc(m.fullName(goFunc), retType, params...)
	f.Sig.Variadic = goFunc.Signature.Variadic()
	// Index LLVM IR function declaration.
	m.globals[goFunc] = f

	// Index anonymous functions declared in fn.
	for _, goAnonFunc := range goFunc.AnonFuncs {
		if err := m.indexFunc(goAnonFunc); err != nil {
			return errors.WithStack(err)
		}
	}
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
	dbg.Println("emitFunc")
	// Index Go SSA function parameters (including receiver of methods).
	fn := m.NewFunc(goFunc)
	dbg.Println("   funcName:", fn.Func.Name())
	for i, goParam := range goFunc.Params {
		param := fn.Func.Params[i]
		fn.locals[goParam] = param
	}
	// Index Go SSA basic blocks by creating corresponding LLVM IR basic blocks.
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
	// Process basic blocks in dominator tree preorder.
	// TODO: ensure stable sorting of basic blocks with equal dominance (sort by
	// basic block index).
	done := make(map[*ssa.BasicBlock]bool)
	goBlocks := goFunc.DomPreorder()
	for len(done) < len(goBlocks) {
		prev := len(done)
	loop:
		for _, goBlock := range goBlocks {
			if done[goBlock] {
				continue
			}
			// Before emitting block containing phi instruction, ensure that all
			// predecessors of goBlock have been emitted.
			if containsPhi(goBlock) {
				for _, goPred := range goBlock.Preds {
					if !done[goPred] {
						continue loop
					}
				}
			}
			done[goBlock] = true
			if err := fn.emitBlock(goBlock); err != nil {
				return errors.WithStack(err)
			}
		}
		if prev == len(done) {
			goFunc.WriteTo(os.Stderr)
			var remaining []string
			for _, goBlock := range goBlocks {
				if !done[goBlock] {
					blockName := getBlockName(goBlock.Index)
					remaining = append(remaining, blockName)
				}
				sort.Strings(remaining)
			}
			warn.Printf("remaining basic blocks with cyclic references: %v", strings.Join(remaining, ", "))
			panic(fmt.Errorf("unable to process basic blocks of %q; cyclic predecessor dependency detected", m.fullName(goFunc)))
		}
	}

	// Compile anonymous functions declared in fn.
	for _, goAnonFunc := range goFunc.AnonFuncs {
		if err := m.emitFunc(goAnonFunc); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// ### [ Helper functions ] ####################################################

// containsPhi reports whether the given basic block contains a phi instruction.
func containsPhi(goBlock *ssa.BasicBlock) bool {
	for _, goInst := range goBlock.Instrs {
		if _, ok := goInst.(*ssa.Phi); ok {
			return true
		}
	}
	return false
}
