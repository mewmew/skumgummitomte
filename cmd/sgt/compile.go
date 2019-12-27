package main

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

// compileMember compiles the given Go SSA member into LLVM IR.
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
