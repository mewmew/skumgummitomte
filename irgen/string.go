package irgen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
)

// emitStringLit compiles the given Go string literal into LLVM IR, emitting to
// m. An LLVM IR global variable holding the contents of the Go string literal
// is created if not already present.
func (m *Module) emitStringLit(s string) *ir.Global {
	m.stringsMutex.Lock()
	defer m.stringsMutex.Unlock()
	if g, ok := m.strings[s]; ok {
		return g
	}
	// Create new LLVM IR global variable holding the contents of the Go string
	// literal.
	strLit := irconstant.NewCharArrayFromString(s)
	strName := m.nextStrName()
	g := m.Module.NewGlobalDef(strName, strLit)
	m.strings[s] = g
	return g
}

// nextStrName returns the next available global variable name to assign an LLVM
// IR global variable holding the contents of a Go string literal.
//
// Note: only thread-safe if invoked thorugh emitStringLit.
func (m *Module) nextStrName() string {
	strName := fmt.Sprintf("str_%04d", m.curStrNum)
	m.curStrNum++
	return strName
}
