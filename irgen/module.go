package irgen

import (
	"sync"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"
	"golang.org/x/tools/go/ssa"
)

// Module is an LLVM IR module generator.
type Module struct {
	// Output LLVM IR module.
	*ir.Module
	// Input Go SSA package.
	goPkg *ssa.Package

	// Maps from global Go SSA value to corresponding LLVM IR value in the LLVM
	// IR module being generated.
	globals map[ssa.Value]irvalue.Value
	// Map from predeclared Go type name to LLVM IR type.
	predeclaredTypes map[string]irtypes.Type
	// Map from predeclared Go function name to LLVM IR function.
	predeclaredFuncs map[string]*ir.Func

	// Mutex to ensure that access to strings and curStrNum is thread-safe.
	stringsMutex sync.Mutex
	// Map from Go string literal to LLVM IR global variable holding the LLVM IR
	// character array constant of the given Go string literal.
	strings map[string]*ir.Global
	// Current string literal number to be used when assigning unique names to
	// LLVM IR global variables holding LLVM IR character array constants of Go
	// string literals.
	curStrNum int
}

// NewModule return a new LLVM IR module generator for the given Go SSA package.
func NewModule(goPkg *ssa.Package) *Module {
	return &Module{
		Module:           ir.NewModule(),
		goPkg:            goPkg,
		globals:          make(map[ssa.Value]irvalue.Value),
		predeclaredTypes: make(map[string]irtypes.Type),
		predeclaredFuncs: make(map[string]*ir.Func),
		strings:          make(map[string]*ir.Global),
	}
}
