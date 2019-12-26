package main

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	"golang.org/x/tools/go/ssa"
)

// generator is an LLVM IR generator for the given SSA Go package.
type generator struct {
	// Input Go package.
	pkg *ssa.Package
	// Output LLVM IR module.
	module *ir.Module

	// Map from constant name to LLVM IR constant.
	consts map[string]constant.Constant
	// Map from global name to LLVM IR global.
	globals map[string]*ir.Global
	// Map from function name to LLVM IR function.
	funcs map[string]*ir.Func
	// Map from type name to LLVM IR type.
	types map[string]irtypes.Type
	// Map from predeclared type name to LLVM IR type.
	predeclaredTypes map[string]irtypes.Type
}

// newGenerator returns a new LLVM IR generator for the given SSA Go package.
func newGenerator(pkg *ssa.Package) *generator {
	return &generator{
		pkg:              pkg,
		module:           ir.NewModule(),
		consts:           make(map[string]constant.Constant),
		globals:          make(map[string]*ir.Global),
		funcs:            make(map[string]*ir.Func),
		types:            make(map[string]irtypes.Type),
		predeclaredTypes: make(map[string]irtypes.Type),
	}
}
