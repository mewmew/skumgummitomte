package main

import (
	"sync"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	"golang.org/x/tools/go/ssa"
)

// generator is an LLVM IR generator for the given SSA Go package.
type generator struct {
	// Input Go package.
	pkg *ssa.Package
	// Output LLVM IR module.
	module *ir.Module

	// Map from Go constant name to LLVM IR constant.
	consts map[string]irconstant.Constant
	// Map from Go global name to LLVM IR global.
	globals map[string]*ir.Global
	// Map from Go function name to LLVM IR function.
	funcs map[string]*ir.Func
	// Map from predeclared Go functions name to LLVM IR functions.
	predeclaredFuncs map[string]*ir.Func
	// Map from Go type name to LLVM IR type.
	types map[string]irtypes.Type
	// Map from predeclared Go type name to LLVM IR type.
	predeclaredTypes map[string]irtypes.Type

	// Mutex to ensure that access to strings and curStrNum is thread-safe.
	stringsMutex sync.Mutex
	// Map from Go string literal to LLVM IR global variable holding the LLVM IR
	// character array constant of the given Go string literal.
	strings map[string]*ir.Global
	// Current string literal number to be used when assigning unique names to
	// LLVM IR global variables holding LLVM IR character array constants of Go
	// string literals.
	curStrNum int
}

// newGenerator returns a new LLVM IR generator for the given SSA Go package.
func newGenerator(pkg *ssa.Package) *generator {
	return &generator{
		pkg:              pkg,
		module:           ir.NewModule(),
		consts:           make(map[string]irconstant.Constant),
		globals:          make(map[string]*ir.Global),
		funcs:            make(map[string]*ir.Func),
		predeclaredFuncs: make(map[string]*ir.Func),
		types:            make(map[string]irtypes.Type),
		predeclaredTypes: make(map[string]irtypes.Type),
		strings:          make(map[string]*ir.Global),
	}
}
