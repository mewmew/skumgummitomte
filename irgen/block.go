package irgen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa"
)

// --- [ get ] -----------------------------------------------------------------

// getBlock returns the LLVM IR basic block corresponding to the given Go SSA
// basic block.
//
// Pre-condition: index basic blocks of fn.
func (fn *Func) getBlock(goBlock *ssa.BasicBlock) *ir.Block {
	// Lookup indexed LLVM IR basic block of Go SSA basic block.
	block, ok := fn.blocks[goBlock]
	if !ok {
		// Pre-condition invalidated, basic block not indexed. This is a fatal
		// error and indicates a bug in irgen.
		panic(fmt.Errorf("unable to locate indexed LLVM IR basic block of Go SSA basic block with index %d (%q)", goBlock.Index, goBlock.Comment))
	}
	return block
}

// --- [ index ] ---------------------------------------------------------------

// indexBlock indexes the given Go SSA basic block, creating a corresponding
// LLVM IR basic block, emitting to fn.
func (fn *Func) indexBlock(goBlock *ssa.BasicBlock) error {
	// Generate LLVM IR basic block.
	blockName := getBlockName(goBlock.Index)
	block := ir.NewBlock(blockName)
	// Index LLVM IR basic block.
	fn.blocks[goBlock] = block
	return nil
}

// --- [ compile ] -------------------------------------------------------------

// TODO: check if the pre-condition of emitBlock should also include "index
// local values of fn".

// emitBlock compiles the given Go SSA basic block into LLVM IR, emitting to fn.
//
// Pre-condition: index basic blocks of fn.
func (fn *Func) emitBlock(goBlock *ssa.BasicBlock) error {
	block := fn.getBlock(goBlock)
	fn.Func.Blocks = append(fn.Func.Blocks, block)
	fn.cur = block
	for _, goInst := range goBlock.Instrs {
		if err := fn.emitInst(goInst); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// ### [ Helper functions ] ####################################################

// getBlockName returns the LLVM IR basic block name based on the given Go SSA
// basic block index.
func getBlockName(index int) string {
	return fmt.Sprintf("block_%04d", index)
}
