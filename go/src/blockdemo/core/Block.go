package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index        int64 //区块编号
	Timestamp    int64
	PreBlockHash string
	Hash         string

	Data string
}

//计算hash值
func calculateHash(b Block) string {
	blockData := string(b.Index) + string(b.Timestamp) + b.PreBlockHash + b.Data
	hashInBytes := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hashInBytes[:])
}

//生成新区块
func GenerateNewBlock(preBlock Block, data string) Block {
	newBlock := Block{}
	newBlock.Index = preBlock.Index + 1
	newBlock.Timestamp = time.Now().Unix()
	newBlock.PreBlockHash = preBlock.Hash
	newBlock.Data = data
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

//创世区块
func generateGenesisBlock() Block {
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.Hash = ""
	return GenerateNewBlock(preBlock, "Genesis Block")

}
