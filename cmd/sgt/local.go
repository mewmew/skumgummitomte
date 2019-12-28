package main

import (
	"fmt"

	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// getLocal returns the LLVM IR value based on the local variable name of the
// corresponding Go SSA local variable.
func (fn *Func) getLocal(localName string) value.Value {
	local, ok := fn.locals[localName]
	if !ok {
		panic(fmt.Errorf("unable to locate LLVM IR value corresponding to Go local variable with name %q", localName))
	}
	return local
}

// defLocal stores v to the local variable of the given name.
func (fn *Func) defLocal(localName string, v value.Value) {
	elemType := v.Type()
	local, ok := fn.locals[localName]
	if !ok {
		local = fn.Entry.NewAlloca(elemType)
		local.SetName(localName)
		fn.locals[localName] = local
	}
	fn.Cur.NewStore(v, local)
}

// useLocal loads and returns the value of the local variable of the given name.
func (fn *Func) useLocal(localName string) value.Value {
	local := fn.getLocal(localName)
	elemType := local.Type().(*irtypes.PointerType).ElemType
	return fn.Cur.NewLoad(elemType, local)
}
