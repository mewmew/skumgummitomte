// Package irgen translates Go SSA code to LLVM IR.
package irgen

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/llir/llvm/ir"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var (
	// dbg is a logger with the "irgen:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("irgen:")+" ", 0)
	// warn is a logger with the "irgen:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("irgen:")+" ", 0)
)

// Default output writer for ssa debug messages.
var ssaDebugWriter io.Writer = os.Stderr

// SetDebugOutput sets the output writer for debug messages to w.
func SetDebugOutput(w io.Writer) {
	dbg.SetOutput(w)
	ssaDebugWriter = w
}

// SetWarningOutput sets the output writer for warning messages to w.
func SetWarningOutput(w io.Writer) {
	warn.SetOutput(w)
}

// CompilePackage compiles the given Go SSA package into an LLVM IR module.
func CompilePackage(goPkg *ssa.Package) (*ir.Module, error) {
	dbg.Println("CompilePackage")
	dbg.Println("   goPkg:", goPkg.Pkg.Name())
	// TODO: remove debug output.
	goPkg.WriteTo(ssaDebugWriter)
	// Create LLVM IR module generator for the given Go SSA package.
	m := NewModule(goPkg)
	// Initialize LLVM IR types corresponding to the predeclared Go types.
	m.initPredeclaredTypes()
	// Initialize LLVM IR functions corresponding to the predeclared Go
	// functions.
	m.initPredeclaredFuncs()
	// Sort type definitions of Go SSA package.
	var goTypeDefs []*ssa.Type
	for _, goMember := range goPkg.Members {
		goTypeDef, ok := goMember.(*ssa.Type)
		if !ok {
			continue
		}
		goTypeDefs = append(goTypeDefs, goTypeDef)
	}
	sort.Slice(goTypeDefs, func(i, j int) bool {
		return goTypeDefs[i].RelString(nil) < goTypeDefs[j].RelString(nil)
	})
	// Index type definitions of Go SSA package.
	// TODO: index type definitions.
	// Compile type definitions of Go SSA package.
	for _, goTypeDef := range goTypeDefs {
		if err := m.emitType(goTypeDef); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Sort functions of Go SSA package.
	goFuncsMap := ssautil.AllFunctions(goPkg.Prog)
	var goFuncs []*ssa.Function
	for goFunc := range goFuncsMap {
		goFuncs = append(goFuncs, goFunc)
	}
	sort.Slice(goFuncs, func(i, j int) bool {
		return m.fullName(goFuncs[i]) < m.fullName(goFuncs[j])
	})
	// Index functions of Go SSA package.
	for _, goFunc := range goFuncs {
		if err := m.indexFunc(goFunc); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	// Sort member names of Go SSA package.
	memberNames := make([]string, 0, len(goPkg.Members))
	for memberName := range goPkg.Members {
		memberNames = append(memberNames, memberName)
	}
	sort.Strings(memberNames)
	// Index SSA members of Go SSA package.
	for _, memberName := range memberNames {
		goMember := goPkg.Members[memberName]
		// TODO: skip function? already indexed above.
		if err := m.indexMember(goMember); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	// Compile functions of Go SSA package.
	for _, goFunc := range goFuncs {
		// Only compile function body of definitions in goPkg.
		if goFunc.Pkg != goPkg {
			continue
		}
		if err := m.emitFunc(goFunc); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	// Compile SSA members of Go SSA package.
	for _, memberName := range memberNames {
		goMember := goPkg.Members[memberName]
		// TODO: skip function? already indexed above.
		if err := m.emitMember(goMember); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return m.Module, nil
}

// --- [ index ] ---------------------------------------------------------------

// indexMember indexes the given Go SSA member, creating a corresponding LLVM IR
// construct, emitting to m.
func (m *Module) indexMember(goMember ssa.Member) error {
	switch goMember := goMember.(type) {
	case *ssa.NamedConst:
		// TODO: evaluate if we need to index named constants or not before
		// resolving them. What cycles may exist?
		//return m.indexNamedConst(goMember)
		return nil // ignore indexing *ssa.NamedConst for now
	case *ssa.Global:
		return m.indexGlobal(goMember)
	case *ssa.Function:
		// handled by indexFunc explicitly.
		return nil
		//return m.indexFunc(goMember)
	case *ssa.Type:
		// TODO: evaluate if we need to index type definitions or not before
		// resolving them. What cycles may exist?
		//return m.indexType(goMember)
		return nil // ignore indexing *ssa.Type for now
	default:
		panic(fmt.Errorf("support for SSA member %T (%q) not yet implemented", goMember, goMember.Name()))
	}
}

// --- [ compile ] -------------------------------------------------------------

// emitMember compiles the given Go SSA member into LLVM IR, emitting to m.
//
// Pre-condition: index global members of m.
func (m *Module) emitMember(goMember ssa.Member) error {
	switch goMember := goMember.(type) {
	// TODO: uncomment.
	case *ssa.NamedConst:
		return m.emitNamedConst(goMember)
	case *ssa.Global:
		return m.emitGlobal(goMember)
	case *ssa.Function:
		// handled by emitFunc explicitly.
		return nil
		//return m.emitFunc(goMember)
	case *ssa.Type:
		// handled by emitType explicitly.
		return nil
		//return m.emitType(goMember)
	default:
		panic(fmt.Errorf("support for SSA member %T (%q) not yet implemented", goMember, goMember.Name()))
	}
}
