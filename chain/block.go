package chain

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Block in the blockchain
type Block struct {
	index        uint64    // 8 bytes
	timestamp    time.Time // 8 bytes for unix
	previousHash [32]byte  // 32 bytes
	hash         [32]byte  // 32 bytes
	data         string    // len(data)

	difficulty int
	nonce      uint64
	noncedHash [32]byte
}

func (block *Block) toByteArray() []byte {
	log.Debug("Converting block to byte array")
	ba := make([]byte, 8+8+32+len(block.data))

	binary.LittleEndian.PutUint64(ba, block.index)
	binary.LittleEndian.PutUint64(ba[8:], uint64(block.timestamp.Unix()))

	copy(ba[16:], block.previousHash[:])
	copy(ba[48:], block.data)
	return ba
}

func (block *Block) String() string {
	return fmt.Sprintf(`
******************************** BLOCK *********************************************
* Index: %d
* Timestamp: %s
* Previous Hash: %x
* Base Hash: %x
* Data: %s
* Difficulty: %d
* Nonce: %d
* Nonced Hash: %x
******************************** ENDBLOCK ******************************************
`,
		block.index,
		block.timestamp.Format(time.ANSIC),
		block.previousHash,
		block.hash,
		block.data,
		block.difficulty,
		block.nonce,
		block.noncedHash)
}

func (block *Block) calculateHash() [32]byte {
	log.Debug("Calculating hash")
	hash := sha256.Sum256(block.toByteArray())
	return hash
}

// CheckNoncedHash checks the hash for a difficulty match
func CheckNoncedHash(hash [32]byte, difficulty int) bool {
	numZeroBytes := (difficulty / 8)
	expectedBytes := make([]byte, numZeroBytes)
	for i := range expectedBytes {
		if expectedBytes[i] != hash[i] {
			return false
		}
	}
	lastByte := int(hash[numZeroBytes+1])
	for i := 0; i < difficulty%8; i++ {
		bitMask := 1 << uint(i%8)
		if lastByte&bitMask > 0 {
			return false
		}
	}
	return true
}

// NoncedHash calculates nonced hash of a block
func (block *Block) NoncedHash() [32]byte {
	withNonce := make([]byte, 32+8)
	copy(withNonce[:], block.hash[:])
	binary.LittleEndian.PutUint64(withNonce[32:], block.nonce)
	noncedHash := sha256.Sum256(withNonce)
	log.Debug(fmt.Sprintf("Calculated nonced hash %x", noncedHash))
	return noncedHash
}

// LoadHash calculates the hash of a block
func (block *Block) LoadHash() [32]byte {
	block.hash = block.calculateHash()
	return block.hash
}

// GenesisBlock generates a beginning block for the chain
func GenesisBlock(data string) *Block {
	b := new(Block)
	b.index = 0
	b.timestamp = time.Now()
	b.data = data
	b.nonce = 0
	b.difficulty = 0
	b.LoadHash()
	return b
}

// NewBlock constructs a new block from the previous
func NewBlock(previous *Block, data string, difficulty int) *Block {
	b := new(Block)
	b.index = previous.index + 1
	b.timestamp = time.Now()
	b.previousHash = previous.noncedHash
	b.data = data
	b.difficulty = difficulty
	b.LoadHash()
	return b
}

// IsValidNewBlock checks if the new block is valid
func IsValidNewBlock(previous *Block, new *Block) bool {
	log.Debug(fmt.Sprintf("Checking new block %d with hash %x", new.index, new.hash))
	if previous.index+1 != new.index {
		log.Warn("Index mismatch")
		return false
	}
	if previous.noncedHash != new.previousHash {
		log.Warn("Previous hash does not match new hash")
		return false
	}
	if new.calculateHash() != new.hash {
		log.Warn("Calculated hash does not match stored hash")
		return false
	}
	if !CheckNoncedHash(new.noncedHash, new.difficulty) {
		log.Warn("Difficulty does not pass tests.")
		return false
	}
	return true
}

// GenerateNonce generates a nonce for a block
func (block *Block) GenerateNonce() {
	for {
		noncedHash := block.NoncedHash()
		if CheckNoncedHash(noncedHash, block.difficulty) {
			block.noncedHash = noncedHash
			return
		}
		block.nonce++
	}
}
