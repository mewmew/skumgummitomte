package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
)

// globalFromStringLit returns the LLVM IR global variable holding the LLVM IR
// character array constant of the given Go string literal. A new LLVM IR global
// variable is created for the Go string literal if not already present in the
// LLVM IR module.
func (gen *generator) globalFromStringLit(s string) *ir.Global {
	gen.stringsMutex.Lock()
	defer gen.stringsMutex.Unlock()
	if g, ok := gen.strings[s]; ok {
		return g
	}
	return gen.addStringLit(s)
}

// addStringLit adds a new global variable to the LLVM IR module being generated
// holding the LLVM IR character array constant of the given Go string literal.
//
// Note: only thread-safe if invoked thorugh globalFromStringLit.
func (gen *generator) addStringLit(s string) *ir.Global {
	strLit := irconstant.NewCharArrayFromString(s)
	strName := gen.nextStrName()
	g := gen.module.NewGlobalDef(strName, strLit)
	gen.strings[s] = g
	return g
}

// nextStrName returns the next available global variable name to assign an LLVM
// IR global variable holding the LLVM IR character array constant of a Go
// string literal.
//
// Note: only thread-safe if invoked thorugh globalFromStringLit.
func (gen *generator) nextStrName() string {
	strName := fmt.Sprintf("str_%04d", gen.curStrNum)
	gen.curStrNum++
	return strName
}
