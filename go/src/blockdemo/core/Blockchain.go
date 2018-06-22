package core

import (
	"fmt"
	"log"
)

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AppendBlock(newBlock *Block) {
	if len(bc.Blocks)==0 {
		bc.Blocks=append(bc.Blocks,newBlock)
		return
	}
	if isValidBlock(*newBlock, *bc.Blocks[len(bc.Blocks)-1]) {
		bc.Blocks = append(bc.Blocks, newBlock)
	} else {
		log.Fatal("invalid block")
	}
}

func NewBlockchain() *Blockchain {
	genesisBlock := generateGenesisBlock()
	blockChain := Blockchain{}
	blockChain.AppendBlock(&genesisBlock)
	return &blockChain
}

func (bc *Blockchain) SendData(data string) {
	preBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := GenerateNewBlock(*preBlock, data)
	bc.AppendBlock(&newBlock)
}

//验证区块合法性
func isValidBlock(newBlock Block, oldBlock Block) bool {
	if newBlock.Index-1 != oldBlock.Index {
		return false
	}
	if newBlock.PreBlockHash != oldBlock.Hash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func (bc *Blockchain) Print() {
	for _, block := range bc.Blocks {
		fmt.Println("Index : %d ,Prev.Hash: %s ,Curr.Hash: %s,Data: %s,Timestamp:%d\n", block.Index, block.PreBlockHash, block.Hash, block.Data, block.Timestamp)
		//fmt.Println("Index",block.Index)
		//fmt.Println("timeStamp",block.Timestamp)
	}
}
