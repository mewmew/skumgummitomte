package irgen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	irenum "github.com/llir/llvm/ir/enum"
	"golang.org/x/tools/go/ssa"
)

// --- [ get ] -----------------------------------------------------------------

// getGlobal returns the LLVM IR global corresponding to the given Go SSA
// global.
//
// Pre-condition: index globals of m.
func (m *Module) getGlobal(goGlobal *ssa.Global) *ir.Global {
	// Lookup indexed LLVM IR global of Go SSA global.
	global, ok := m.globals[goGlobal]
	if !ok {
		// Pre-condition invalidated, global not indexed. This is a fatal error
		// and indicates a bug in irgen.
		panic(fmt.Errorf("unable to locate indexed LLVM IR global of Go SSA global %q", m.fullName(goGlobal)))
	}
	return global.(*ir.Global)
}

// --- [ index ] ---------------------------------------------------------------

// indexGlobal indexes the given Go SSA global, creating a corresponding LLVM IR
// global variable, emitting to m. The external boolean indicates whether the
// global variable is defined in an external Go package.
func (m *Module) indexGlobal(goGlobal *ssa.Global, external bool) error {
	// Generate LLVM IR global variable declaration, emitting to m.
	globalType := m.irTypeFromGo(goGlobal.Type())
	global := m.Module.NewGlobal(m.fullName(goGlobal), globalType)
	// Add external linkage to global variable defined in external Go package.
	if external {
		global.Linkage = irenum.LinkageExternal
	}
	// Index LLVM IR global variable declaration.
	m.globals[goGlobal] = global
	return nil
}

// --- [ compile ] -------------------------------------------------------------

// emitGlobal compiles the given Go SSA global into LLVM IR, emitting to m.
//
// Pre-condition: index globals of m.
func (m *Module) emitGlobal(goGlobal *ssa.Global) error {
	global := m.getGlobal(goGlobal)
	// TODO: check how initialization of Go globals work. Is it always done
	// through an `init` function? If so, using a zeroinitializer is always
	// correct. If not, we should initialize the LLVM IR value using the Go
	// global initializer.
	global.Init = irconstant.NewZeroInitializer(global.ContentType)
	return nil
}
