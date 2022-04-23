package Blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	Transaction string
	Reward      string
	PrevPointer *Block
	Hash        [32]byte
	PrevHash    [32]byte
}

func DeriveHash(transaction string) [32]byte {
	return sha256.Sum256([]byte(transaction))
}

func InsertBlock(Transaction string, reward string, chainHead *Block) *Block {
	if chainHead == nil {
		return &Block{Transaction, reward, nil, DeriveHash(Transaction), [32]byte{}}
	}
	return &Block{Transaction, reward, chainHead, DeriveHash(Transaction), DeriveHash(chainHead.Transaction)}
}

func ListBlocks(block *Block) {
	fmt.Println("Blockchain")
	fmt.Printf("%-10s\t", "Transaction")
	fmt.Printf("%-10s\n", "Reward")
	for currBlock := block; currBlock != nil; currBlock = currBlock.PrevPointer {
		fmt.Printf("%-10s\t", currBlock.Transaction)
		fmt.Printf("%-10s\n", currBlock.Reward)
		//fmt.Printf("Hash: %x\n", currBlock.Hash)
		//fmt.Printf("Previous Hash: %x\n", currBlock.PrevHash)
	}
}

func ReturnBlocks(block *Block) string {
	var blockchain string
	for currBlock := block; currBlock != nil; currBlock = currBlock.PrevPointer {
		blockchain += "\nTransaction: " + currBlock.Transaction + "\n"
		blockchain += "\nReward: " + currBlock.Reward + "\n"
		hash := hex.EncodeToString([]byte(currBlock.Hash[:]))
		blockchain += "Hash: " + hash + "\n"
		PrevHash := hex.EncodeToString([]byte(currBlock.PrevHash[:]))
		blockchain += "Previous Hash: " + PrevHash + "\n"
	}

	return blockchain
}

func CountBlocks(block *Block) int {
	count := 0
	for p := block; p != nil; p = p.PrevPointer {
		count += 1
	}
	return count
}

func VerifyChain(block *Block) bool {
	for p := block; p != nil; p = p.PrevPointer {
		if p.PrevPointer != nil {
			if p.PrevHash != p.PrevPointer.Hash {
				return false
			}
		}
	}
	return true
}
