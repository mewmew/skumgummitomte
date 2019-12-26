package main

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

// indexMember indexes the given SSA member into LLVM IR.
func (gen *generator) indexMember(memberName string, member ssa.Member) error {
	switch member := member.(type) {
	case *ssa.NamedConst:
		return gen.indexNamedConst(memberName, member)
	case *ssa.Global:
		return gen.indexGlobal(memberName, member)
	case *ssa.Function:
		return gen.indexFunction(memberName, member)
	case *ssa.Type:
		return gen.indexType(memberName, member)
	default:
		panic(fmt.Errorf("support for SSA member %T not yet implemented", member))
	}
}

// indexNamedConst indexes the given SSA NamedConst into LLVM IR.
func (gen *generator) indexNamedConst(constName string, c *ssa.NamedConst) error {
	fmt.Println("indexNamedConst")
	// TODO: remove debug output.
	fmt.Println(c)
	fmt.Println()
	return nil
}

// indexGlobal indexes the given SSA Global into LLVM IR.
func (gen *generator) indexGlobal(globalName string, global *ssa.Global) error {
	fmt.Println("indexGlobal")
	// TODO: remove debug output.
	fmt.Println(global)
	fmt.Println()
	return nil
}

// indexFunction indexes the given SSA Function into LLVM IR.
func (gen *generator) indexFunction(funcName string, f *ssa.Function) error {
	fmt.Println("indexFunction")
	// TODO: remove debug output.
	fmt.Println(f)
	fmt.Println()
	return nil
}

// indexType indexes the given SSA Type into LLVM IR.
func (gen *generator) indexType(typeName string, typ *ssa.Type) error {
	fmt.Println("indexType")
	// TODO: remove debug output.
	fmt.Println(typ)
	fmt.Println()
	return nil
}
