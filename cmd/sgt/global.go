package main

import (
	"fmt"

	irconstant "github.com/llir/llvm/ir/constant"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// compileGlobal compiles the given Go SSA global into LLVM IR.
func (gen *generator) compileGlobal(globalName string, goGlobal *ssa.Global) error {
	// TODO: remove debug output.
	fmt.Println("compileGlobal")
	fmt.Println(goGlobal)
	fmt.Println()

	// Locate output LLVM IR global.
	global, ok := gen.globals[globalName]
	if !ok {
		return errors.Errorf("unable to locate LLVM IR global %q", globalName)
	}
	_ = global

	// TODO: remove once support for global initializer is added.
	global.Init = irconstant.NewZeroInitializer(global.Typ.ElemType)

	return nil
}
