package main

import (
	"flag"
	"log"
	"os"
	"sort"
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

// sgt compiles the Go packages specified by package path patterns into LLVM IR
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

// compilePackage compiles the given Go package into an LLVM IR module.
func compilePackage(pkg *ssa.Package) error {
	gen := newGenerator(pkg)
	// TODO: remove debug output.
	pkg.WriteTo(os.Stdout)
	// Sort member names.
	memberNames := make([]string, 0, len(pkg.Members))
	for memberName := range pkg.Members {
		memberNames = append(memberNames, memberName)
	}
	sort.Strings(memberNames)
	// Index SSA members of Go package.
	for _, memberName := range memberNames {
		member := pkg.Members[memberName]
		if err := gen.indexMember(memberName, member); err != nil {
			return errors.WithStack(err)
		}
	}
	// Compile SSA members of Go package.
	for _, memberName := range memberNames {
		member := pkg.Members[memberName]
		if err := gen.compileMember(memberName, member); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
