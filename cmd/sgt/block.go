package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
)

// getBlockName returns the LLVM IR basic block name based on the given basic
// block index of the corresponding Go SSA basic block.
func getBlockName(index int) string {
	return fmt.Sprintf("block_%d", index)
}

// getBlock returns the LLVM IR basic block based on the given basic block index
// of the corresponding Go SSA basic block.
func (fn *Func) getBlock(index int) *ir.Block {
	blockName := getBlockName(index)
	block, ok := fn.blocks[blockName]
	if !ok {
		panic(fmt.Errorf("unable to locate LLVM IR basic block corresponding to Go SSA basic block with index %d", index))
	}
	return block
}
