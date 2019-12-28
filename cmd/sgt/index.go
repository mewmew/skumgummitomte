package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	"golang.org/x/tools/go/ssa"
)

// indexMember indexes the given Go SSA member, creating a scaffolding for the
// corresponding LLVM IR.
func (gen *generator) indexMember(memberName string, member ssa.Member) error {
	switch member := member.(type) {
	case *ssa.NamedConst:
		return gen.indexNamedConst(memberName, member)
	case *ssa.Global:
		return gen.indexGlobal(memberName, member)
	case *ssa.Function:
		return gen.indexFunction(memberName, member)
	case *ssa.Type:
		return gen.indexType(memberName, member)
	default:
		panic(fmt.Errorf("support for SSA member %T not yet implemented", member))
	}
}

// indexNamedConst indexes the given Go SSA named constant, creating a
// scaffolding for the corresponding LLVM IR constant.
func (gen *generator) indexNamedConst(constName string, goConst *ssa.NamedConst) error {
	// TODO: remove debug output.
	dbg.Println("indexNamedConst")
	dbg.Println(goConst)
	dbg.Println()

	panic("not yet implemented")
	return nil
}

// indexGlobal indexes the given Go SSA global, creating a scaffolding for the
// corresponding LLVM IR global variable.
func (gen *generator) indexGlobal(globalName string, goGlobal *ssa.Global) error {
	// TODO: remove debug output.
	dbg.Println("indexGlobal")
	dbg.Println(goGlobal)

	// Generate LLVM IR global variable declaration.
	goType := goGlobal.Type()
	globalType := gen.llTypeFromGoType(goType)
	global := gen.module.NewGlobal(globalName, globalType)
	gen.globals[globalName] = global

	// TODO: remove debug output.
	dbg.Println("global:", global.LLString())
	dbg.Println()
	return nil
}

// indexFunction indexes the given Go SSA function, creating a scaffolding for
// the corresponding LLVM IR function.
func (gen *generator) indexFunction(funcName string, goFunc *ssa.Function) error {
	// TODO: remove debug output.
	dbg.Println("indexFunction")
	dbg.Println(goFunc)

	// TODO: add support for receiver of methods.

	// Convert Go function parameters to equivalent LLVM IR function parameters.
	var params []*ir.Param
	goParams := goFunc.Signature.Params()
	for i := 0; i < goParams.Len(); i++ {
		goParam := goParams.At(i)
		paramName := goParam.Name()
		goParamType := goParam.Type()
		paramType := gen.llTypeFromGoType(goParamType)
		param := ir.NewParam(paramName, paramType)
		params = append(params, param)
	}

	// Convert Go function return types to equivalent LLVM IR function return
	// types.
	var resultTypes []irtypes.Type
	goResults := goFunc.Signature.Results()
	for i := 0; i < goResults.Len(); i++ {
		goResult := goResults.At(i)
		resultName := goResult.Name()
		// TODO: add resultName as field name of (custom) result structure type.
		_ = resultName
		goResultType := goResult.Type()
		resultType := gen.llTypeFromGoType(goResultType)
		resultTypes = append(resultTypes, resultType)
	}
	// Convert multiple return types a single return type by creating a structure
	// type with one field per return type.
	var retType irtypes.Type
	switch len(resultTypes) {
	// void return.
	case 0:
		retType = irtypes.Void
	// single return type.
	case 1:
		retType = resultTypes[0]
	// multiple return types.
	default:
		retType = irtypes.NewStruct(resultTypes...)
	}

	// Generate LLVM IR function declaration.
	f := gen.module.NewFunc(funcName, retType, params...)
	f.Sig.Variadic = goFunc.Signature.Variadic()
	gen.funcs[funcName] = f

	// TODO: remove debug output.
	dbg.Println("f:", f.LLString())
	dbg.Println()
	return nil
}

// indexType indexes the given Go SSA type, creating a scaffolding for the
// corresponding LLVM IR type.
func (gen *generator) indexType(typeName string, goType *ssa.Type) error {
	// TODO: remove debug output.
	dbg.Println("indexType")
	dbg.Println(goType)
	dbg.Println()

	panic("not yet implemented")
	return nil
}
