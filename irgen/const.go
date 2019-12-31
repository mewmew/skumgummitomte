package irgen

import "golang.org/x/tools/go/ssa"

// --- [ index ] ---------------------------------------------------------------

// indexNamedConst indexes the given Go SSA named constant, creating a
// corresponding LLVM IR constant, emitting to m.
func (m *Module) indexNamedConst(goConst *ssa.NamedConst) error {
	panic("not yet implemented")
	return nil
}

// --- [ compile ] -------------------------------------------------------------
