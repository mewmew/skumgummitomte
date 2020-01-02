package irgen

import (
	irconstant "github.com/llir/llvm/ir/constant"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// --- [ index ] ---------------------------------------------------------------

// indexNamedConst indexes the given Go SSA named constant, creating a
// corresponding LLVM IR constant, emitting to m.
func (m *Module) indexNamedConst(goConst *ssa.NamedConst) error {
	panic("not yet implemented")
}

// --- [ compile ] -------------------------------------------------------------

// emitNamedConst compiles the given Go SSA named constant into LLVM IR,
// emitting to m.
func (m *Module) emitNamedConst(goConst *ssa.NamedConst) error {
	dbg.Println("emitNamedConst")
	v := m.irValueFromGo(goConst.Value)
	c, ok := v.(irconstant.Constant)
	if !ok {
		return errors.Errorf("unable to convert value %T into LLVM IR constant, as needed by named constant definition", v)
	}
	m.consts[goConst] = c
	return nil
}
