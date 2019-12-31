// Package irgen translates Go SSA code to LLVM IR.
package irgen

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/llir/llvm/ir"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

var (
	// dbg is a logger with the "irgen:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("irgen:")+" ", 0)
	// warn is a logger with the "irgen:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("irgen:")+" ", 0)
)

// CompilePackage compiles the given Go SSA package into an LLVM IR module.
func CompilePackage(goPkg *ssa.Package) (*ir.Module, error) {
	// TODO: remove debug output.
	goPkg.WriteTo(os.Stdout)
	// Create LLVM IR module generator for the given Go SSA package.
	m := NewModule(goPkg)
	// Initialize LLVM IR types corresponding to the predeclared Go types.
	m.initPredeclaredTypes()
	// Initialize LLVM IR functions corresponding to the predeclared Go
	// functions.
	m.initPredeclaredFuncs()
	// Sort member names of SSA Go package.
	memberNames := make([]string, 0, len(goPkg.Members))
	for memberName := range goPkg.Members {
		memberNames = append(memberNames, memberName)
	}
	sort.Strings(memberNames)
	// Index SSA members of Go package.
	for _, memberName := range memberNames {
		member := goPkg.Members[memberName]
		if err := m.indexMember(member); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	// Compile SSA members of Go package.
	for _, memberName := range memberNames {
		member := goPkg.Members[memberName]
		if err := m.emitMember(member); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return m.Module, nil
}

// --- [ index ] ---------------------------------------------------------------

// indexMember indexes the given Go SSA member, creating a corresponding LLVM IR
// construct, emitting to m.
func (m *Module) indexMember(member ssa.Member) error {
	switch member := member.(type) {
	case *ssa.NamedConst:
		return m.indexNamedConst(member)
	case *ssa.Global:
		return m.indexGlobal(member)
	case *ssa.Function:
		return m.indexFunc(member)
	case *ssa.Type:
		return m.indexType(member)
	default:
		panic(fmt.Errorf("support for SSA member %T (%q) not yet implemented", member, member.Name()))
	}
}

// --- [ compile ] -------------------------------------------------------------

// emitMember compiles the given Go SSA member into LLVM IR, emitting to m.
//
// Pre-condition: index global members of m.
func (m *Module) emitMember(member ssa.Member) error {
	switch member := member.(type) {
	// TODO: uncomment.
	//case *ssa.NamedConst:
	//	return m.emitNamedConst(member)
	case *ssa.Global:
		return m.emitGlobal(member)
	case *ssa.Function:
		return m.emitFunc(member)
	// TODO: uncomment.
	//case *ssa.Type:
	//	return m.emitType(member)
	default:
		panic(fmt.Errorf("support for SSA member %T (%q) not yet implemented", member, member.Name()))
	}
}
