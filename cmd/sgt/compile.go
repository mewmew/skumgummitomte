package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
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
	types map[string]types.Type
}

// newGenerator returns a new LLVM IR generator for the given SSA Go package.
func newGenerator(pkg *ssa.Package) *generator {
	return &generator{
		pkg:     pkg,
		module:  ir.NewModule(),
		consts:  make(map[string]constant.Constant),
		globals: make(map[string]*ir.Global),
		funcs:   make(map[string]*ir.Func),
		types:   make(map[string]types.Type),
	}
}

// compileMember compiles the given SSA member into LLVM IR.
func (gen *generator) compileMember(memberName string, member ssa.Member) error {
	switch member := member.(type) {
	case *ssa.NamedConst:
		return gen.compileNamedConst(memberName, member)
	case *ssa.Global:
		return gen.compileGlobal(memberName, member)
	case *ssa.Function:
		return gen.compileFunction(memberName, member)
	case *ssa.Type:
		return gen.compileType(memberName, member)
	default:
		panic(fmt.Errorf("support for SSA member %T not yet implemented", member))
	}
}

// compileNamedConst compiles the given SSA NamedConst into LLVM IR.
func (gen *generator) compileNamedConst(constName string, goConst *ssa.NamedConst) error {
	// TODO: remove debug output.
	fmt.Println("compileNamedConst")
	fmt.Println(goConst)
	fmt.Println()
	return nil
}

// compileGlobal compiles the given SSA Global into LLVM IR.
func (gen *generator) compileGlobal(globalName string, goGlobal *ssa.Global) error {
	// TODO: remove debug output.
	fmt.Println("compileGlobal")
	fmt.Println(goGlobal)
	fmt.Println()
	return nil
}

// compileFunction compiles the given SSA Function into LLVM IR.
func (gen *generator) compileFunction(funcName string, goFunc *ssa.Function) error {
	// TODO: remove debug output.
	fmt.Println("compileFunction")
	fmt.Println(goFunc)
	fmt.Println()
	return nil
}

// compileType compiles the given SSA Type into LLVM IR.
func (gen *generator) compileType(typeName string, goType *ssa.Type) error {
	// TODO: remove debug output.
	fmt.Println("compileType")
	fmt.Println(goType)
	fmt.Println()
	return nil
}
