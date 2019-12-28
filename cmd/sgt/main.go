package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var (
	// dbg is a logger with the "sgt:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.CyanBold("sgt:")+" ", 0)
	// warn is a logger with the "sgt:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("sgt:")+" ", 0)
)

const use = `
Usage:

	sgt [OPTION]... package...

Flags:
`

func usage() {
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	var (
		// Output path of LLVM IR module.
		output string
	)
	flag.StringVar(&output, "o", "", "output path of LLVM IR module (default: standard output)")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	pkgPaths := flag.Args()

	// Write to standard output or output file path if specified by -o flag.
	w := os.Stdout
	if len(output) > 0 {
		f, err := os.Create(output)
		if err != nil {
			log.Fatalf("%+v", errors.WithStack(err))
		}
		defer f.Close()
		w = f
	}

	// Compile packages to LLVM IR modules.
	// TODO: figure out a better way to specify output path, as we want each Go
	// package to be written to a dedicated LLVM IR module. Perhaps specify
	// output directory?
	if err := sgt(w, pkgPaths); err != nil {
		log.Fatalf("%+v", err)
	}
}

// sgt compiles the Go packages specified by package path patterns into LLVM IR
// modules.
func sgt(w io.Writer, pkgPaths []string) error {
	// Parse and type-check Go packages.
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	initial, err := packages.Load(cfg, pkgPaths...)
	if err != nil {
		return errors.WithStack(err)
	}
	// Stop early if there are errors in any of the packages.
	if packages.PrintErrors(initial) > 0 {
		return errors.Errorf("packages contain errors (%s)", strings.Join(pkgPaths, ", "))
	}
	// Create SSA packages of Go packages.
	mode := ssa.PrintPackages | ssa.PrintFunctions | ssa.NaiveForm
	prog, pkgs := ssautil.Packages(initial, mode)
	// Build SSA code for all Go packages.
	prog.Build()
	// Compile Go packages to LLVM IR.
	for _, pkg := range pkgs {
		module, err := compilePackage(pkg)
		if err != nil {
			return errors.WithStack(err)
		}
		dbg.Printf("LLVM IR module of %q:\n", pkg.Pkg.Name())
		if _, err := module.WriteTo(w); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// compilePackage compiles the given Go package into an LLVM IR module.
func compilePackage(pkg *ssa.Package) (*ir.Module, error) {
	// TODO: remove debug output.
	pkg.WriteTo(os.Stdout)

	// Create LLVM IR generator for the given Go package.
	gen := newGenerator(pkg)
	// Initialize LLVM IR types corresponding to the predeclared Go types.
	gen.initPredeclaredTypes()
	// Initialize LLVM IR functions corresponding to the predeclared Go functions.
	gen.initPredeclaredFuncs()

	// Sort member names of SSA Go package.
	memberNames := make([]string, 0, len(pkg.Members))
	for memberName := range pkg.Members {
		memberNames = append(memberNames, memberName)
	}
	sort.Strings(memberNames)

	// Index SSA members of Go package.
	for _, memberName := range memberNames {
		member := pkg.Members[memberName]
		if err := gen.indexMember(memberName, member); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Compile SSA members of Go package.
	for _, memberName := range memberNames {
		member := pkg.Members[memberName]
		if err := gen.compileMember(memberName, member); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return gen.module, nil
}
