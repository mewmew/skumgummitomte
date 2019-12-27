package main

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// compileNamedConst compiles the given Go SSA named constant into LLVM IR.
func (gen *generator) compileNamedConst(constName string, goConst *ssa.NamedConst) error {
	// TODO: remove debug output.
	fmt.Println("compileNamedConst")
	fmt.Println(goConst)
	fmt.Println()

	// Locate output LLVM IR constant.
	c, ok := gen.consts[constName]
	if !ok {
		return errors.Errorf("unable to locate LLVM IR constant %q", constName)
	}
	_ = c
	return nil
}
