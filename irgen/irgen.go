// Package irgen translates Go SSA code to LLVM IR.
package irgen

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
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

	// Compile type definitions of Go SSA package and its dependencies.
	emitted := make(map[*ssa.Package]bool)
	if err := m.emitAllTypeDefs(goPkg, emitted); err != nil {
		return nil, errors.WithStack(err)
	}

	// Index members of Go SSA package and its dependencies.
	indexed := make(map[*ssa.Package]bool)
	if err := m.indexAllMembers(goPkg, indexed); err != nil {
		return nil, errors.WithStack(err)
	}

	// Sort member names of Go SSA package.
	goMembers := make([]ssa.Member, 0, len(goPkg.Members))
	for _, goMember := range goPkg.Members {
		goMembers = append(goMembers, goMember)
	}
	sort.Slice(goMembers, func(i, j int) bool {
		return goMembers[i].RelString(nil) < goMembers[j].RelString(nil)
	})

	// Compile members of Go SSA package.
	for _, goMember := range goMembers {
		if err := m.emitMember(goMember); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Hook up forward declaration (function stubs).
	//
	// ref: https://dave.cheney.net/2019/08/20/go-compiler-intrinsics
	var fs []*ir.Func
	var externalFuncs []*ir.Func // TODO: remove
	for goValue := range m.globals {
		goFunc, ok := goValue.(*ssa.Function)
		if !ok {
			continue
		}
		f := m.getFunc(goFunc)
		if goFunc.Pkg != m.goPkg {
			// skip external declarations.
			externalFuncs = append(externalFuncs, f)
			continue
		}
		fs = append(fs, f)
	}
	sort.Slice(externalFuncs, func(i, j int) bool {
		return externalFuncs[i].Name() < externalFuncs[j].Name()
	})
	sort.Slice(fs, func(i, j int) bool {
		return fs[i].Name() < fs[j].Name()
	})
	for _, externalFunc := range externalFuncs {
		dbg.Println("external function:", externalFunc.Name())
	}
	funcMap := make(map[string]*ir.Func)
	for _, f := range fs {
		funcMap[f.Name()] = f
	}
	for _, f := range fs {
		if len(f.Blocks) == 0 {
			bodyName := strings.ToLower(f.Name())
			bodyFunc, ok := funcMap[bodyName]
			if !ok {
				continue
			}
			if !irtypes.Equal(f.Sig, bodyFunc.Sig) {
				warn.Printf("function signature mismatch between %q (%v) and %q (%v)", f.Name(), f.Sig.LLString(), bodyFunc.Name(), bodyFunc.Sig.LLString())
				continue
			}
			entry := f.NewBlock("entry")
			args := make([]irvalue.Value, 0, len(f.Params))
			for _, param := range f.Params {
				args = append(args, param)
			}
			var result irvalue.Value
			callInst := entry.NewCall(bodyFunc, args...)
			if !irtypes.Equal(f.Sig.RetType, irtypes.Void) {
				result = callInst
			}
			entry.NewRet(result)
			dbg.Printf("function body %q added to stub %q", bodyFunc.Name(), f.Name())
		}
	}
	for _, f := range fs {
		if len(f.Blocks) == 0 {
			dbg.Println("forward declaration:", f.Name())
		} else {
			dbg.Println("function definition:", f.Name())
		}
	}
	return m.Module, nil
}

// --- [ index ] ---------------------------------------------------------------

// indexAllMembers indexes the members of the given Go SSA package and its
// dependencies, creating corresponding LLVM IR constructs, emitting to m.
func (m *Module) indexAllMembers(goPkg *ssa.Package, indexed map[*ssa.Package]bool) error {
	if indexed[goPkg] {
		return nil
	}
	indexed[goPkg] = true
	for _, imp := range goPkg.Pkg.Imports() {
		goImpPkg := goPkg.Prog.Package(imp)
		if err := m.indexAllMembers(goImpPkg, indexed); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := m.indexMembers(goPkg); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// indexMembers indexes the members of the given Go SSA package, creating
// corresponding LLVM IR constructs, emitting to m.
func (m *Module) indexMembers(goPkg *ssa.Package) error {
	// Sort member names of Go SSA package.
	dbg.Println("indexing package:", goPkg.Pkg.Path())
	goMembers := make([]ssa.Member, 0, len(goPkg.Members))
	for _, goMember := range goPkg.Members {
		goMembers = append(goMembers, goMember)
	}
	sort.Slice(goMembers, func(i, j int) bool {
		return goMembers[i].RelString(nil) < goMembers[j].RelString(nil)
	})
	for _, goMember := range goMembers {
		dbg.Println("   index member:", goMember.RelString(nil))
	}
	// Index SSA members of Go SSA package.
	external := m.goPkg != goPkg
	for _, goMember := range goMembers {
		if err := m.indexMember(goMember, external); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// indexMember indexes the given Go SSA member, creating a corresponding LLVM IR
// construct, emitting to m. The external boolean indicates whether the Go SSA
// member is defined in an external Go package.
func (m *Module) indexMember(goMember ssa.Member, external bool) error {
	switch goMember := goMember.(type) {
	case *ssa.NamedConst:
		// TODO: index named constant as LLVM IR constant global variable.
		//return m.indexNamedConst(goMember)
		return nil // ignore indexing *ssa.NamedConst for now
	case *ssa.Global:
		m.indexGlobal(goMember, external)
		return nil
	case *ssa.Function:
		return m.indexFunc(goMember)
	case *ssa.Type:
		// TODO: figure out how to index type definitions.
		//return m.indexType(goMember)
		return nil // ignore indexing *ssa.Type for now
	default:
		panic(fmt.Errorf("support for SSA member %T (%q) not yet implemented", goMember, goMember.Name()))
	}
}

// --- [ compile ] -------------------------------------------------------------

// emitAllTypeDefs compiles the type definitions of the given Go SSA package and
// its dependencies into LLVM IR, emitting to m.
//
// Pre-condition: index type definitions of m.
func (m *Module) emitAllTypeDefs(goPkg *ssa.Package, emitted map[*ssa.Package]bool) error {
	if emitted[goPkg] {
		return nil
	}
	emitted[goPkg] = true
	for _, imp := range goPkg.Pkg.Imports() {
		goImpPkg := goPkg.Prog.Package(imp)
		if err := m.emitAllTypeDefs(goImpPkg, emitted); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := m.emitTypeDefs(goPkg); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// emitTypeDefs compiles the type definitions of the given Go SSA package into
// LLVM IR, emitting to m.
//
// Pre-condition: index type definitions of m.
func (m *Module) emitTypeDefs(goPkg *ssa.Package) error {
	// Sort member names of Go SSA package.
	var goTypes []*ssa.Type
	for _, goMember := range goPkg.Members {
		if goType, ok := goMember.(*ssa.Type); ok {
			goTypes = append(goTypes, goType)
		}
	}
	sort.Slice(goTypes, func(i, j int) bool {
		return goTypes[i].RelString(nil) < goTypes[j].RelString(nil)
	})
	// Emit Go SSA type definition of Go SSA package.
	for _, goType := range goTypes {
		if err := m.emitType(goType); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// --- [ compile ] -------------------------------------------------------------

// emitMember compiles the given Go SSA member into LLVM IR, emitting to m.
//
// Pre-condition: index global members of m.
func (m *Module) emitMember(goMember ssa.Member) error {
	switch goMember := goMember.(type) {
	case *ssa.NamedConst:
		return m.emitNamedConst(goMember)
	case *ssa.Global:
		return m.emitGlobal(goMember)
	case *ssa.Function:
		return m.emitFunc(goMember)
	case *ssa.Type:
		// handled by emitType explicitly.
		return nil
		//return m.emitType(goMember)
	default:
		panic(fmt.Errorf("support for SSA member %T (%q) not yet implemented", goMember, goMember.Name()))
	}
}
