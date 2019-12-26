package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	flag.Parse()
	pkgPaths := flag.Args()
	if err := sgt(pkgPaths); err != nil {
		log.Fatalf("%+v", err)
	}
}

// sgt compiles the Go packages specified by package path patterns to LLVM IR
// modules.
func sgt(pkgPaths []string) error {
	// Parse and type-check Go packages.
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	initial, err := packages.Load(cfg, pkgPaths...)
	if err != nil {
		return errors.WithStack(err)
	}
	// Stop early if there are errors in any of the packages.
	if packages.PrintErrors(initial) > 0 {
		return errors.Errorf("packages contain errors (%s)", strings.Join(pkgPaths, ", "))
	}
	// Create SSA packages of Go packages.
	prog, pkgs := ssautil.Packages(initial, ssa.PrintPackages)
	_ = prog
	// Build SSA code for all Go packages.
	for _, pkg := range pkgs {
		if pkg != nil {
			pkg.Build()
		}
	}
	for _, pkg := range pkgs {
		if err := compilePackage(pkg); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// compilePackage compiles the given Go package to an LLVM IR module.
func compilePackage(pkg *ssa.Package) error {
	// TODO: remove debug output.
	pkg.WriteTo(os.Stdout)
	return nil
}
