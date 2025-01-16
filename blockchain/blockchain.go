package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index        int       `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
	Data         string    `json:"data"`
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"-"`
	Nonce        int       `json:"nonce"`
}

func CalculateHash(block Block) string {
	// Преобразуем в JSON для хеширования (или любой другой способ сериализации)
	record, _ := json.Marshal(block)
	h := sha256.Sum256(record)
	return hex.EncodeToString(h[:])
}

type Blockchain struct {
	Chain []Block
}

func NewBlockchain() *Blockchain {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         "Genesis Block",
		PreviousHash: "",
		Hash:         "",
		Nonce:        0,
	}
	genesisBlock.Hash = CalculateHash(genesisBlock)

	return &Blockchain{
		Chain: []Block{genesisBlock},
	}
}

func (bc *Blockchain) GenerateBlock(data string) Block {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		Nonce:        0, // Для PoW
	}

	// Если используем PoW – добавляем логику майнинга
	// иначе просто рассчитываем хеш
	minedBlock := ProofOfWork(newBlock)

	return minedBlock
}

func ProofOfWork(block Block) Block {
	targetPrefix := "0000" // Сложность: 4 ведущих нуля (пример)
	for {
		hash := CalculateHash(block)
		if hash[:len(targetPrefix)] == targetPrefix {
			block.Hash = hash
			break
		}
		block.Nonce++
	}
	return block
}

func (bc *Blockchain) AddBlock(newBlock Block) {
	// Проверяем корректность
	if bc.isBlockValid(newBlock) {
		bc.Chain = append(bc.Chain, newBlock)
	} else {
		fmt.Println("Block is not valid")
	}
}

func (bc *Blockchain) isBlockValid(block Block) bool {
	prevBlock := bc.Chain[len(bc.Chain)-1]

	// Проверяем индекс
	if block.Index != prevBlock.Index+1 {
		fmt.Println("Index is not valid")
		return false
	}
	// Проверяем PreviousHash
	if block.PreviousHash != prevBlock.Hash {
		fmt.Println("PreviousHash is not valid")
		return false
	}
	// Проверяем хеш
	if CalculateHash(block) != block.Hash {
		fmt.Println("Hash is not valid")
		return false
	}

	return true
}
