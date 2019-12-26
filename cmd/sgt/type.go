package main

import (
	"fmt"
	gotypes "go/types"

	irtypes "github.com/llir/llvm/ir/types"
)

// llTypeFromGoType returns the LLVM IR type corresponding to the given Go type.
func llTypeFromGoType(goType gotypes.Type) irtypes.Type {
	switch goType := goType.(type) {
	default:
		panic(fmt.Errorf("support for Go type %T not yet implemented", goType))
	}
}
