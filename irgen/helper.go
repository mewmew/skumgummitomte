package irgen

import gotypes "go/types"

// ### [ Helper functions ] ####################################################

// RelStringer is the interface that wraps the Go SSA RelString method.
type RelStringer interface {
	// RelString returns the full name of the global, qualified by package name,
	// receiver type, etc.
	//
	// Examples:
	//
	//    "math.IsNaN"                  // a package-level function
	//    "(*bytes.Buffer).Bytes"       // a declared method or a wrapper
	//    "(*bytes.Buffer).Bytes$thunk" // thunk (func wrapping method; receiver is param 0)
	//    "(*bytes.Buffer).Bytes$bound" // bound (func wrapping method; receiver supplied by closure)
	//    "main.main$1"                 // an anonymous function in main
	//    "main.init#1"                 // a declared init function
	//    "main.init"                   // the synthesized package initializer
	RelString(from *gotypes.Package) string
}

// fullName returns the full name of the value, qualified by package name,
// receiver type, etc.
func (m *Module) fullName(v RelStringer) string {
	if m.goPkg.Pkg.Name() == "main" {
		// Fully qualified name if global is imported, otherwise name without
		// package path.
		from := m.goPkg.Pkg
		return v.RelString(from)
	}
	// Fully qualified name (with package path).
	return v.RelString(nil)
}
