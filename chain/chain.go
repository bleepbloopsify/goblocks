package chain

import "bytes"

// Chain to hold blocks
type Chain struct {
	blocks []*Block
}

// Genesis generates an initial chain
func Genesis(data string) *Chain {
	chain := new(Chain)
	chain.blocks = []*Block{GenesisBlock(data)}
	return chain
}

// Length of the chain
func (chain *Chain) Length() int {
	return len(chain.blocks)
}

// AddBlock adds a block to the chain
func (chain *Chain) AddBlock(block *Block) bool {
	if IsValidNewBlock(chain.blocks[len(chain.blocks)-1], block) {
		chain.blocks = append(chain.blocks, block)
		return true
	}
	return false
}

// AddData adds data to the blockchain
func (chain *Chain) AddData(data string, difficulty int) bool {
	block := NewBlock(chain.blocks[len(chain.blocks)-1], data, difficulty)
	block.GenerateNonce()
	return chain.AddBlock(block)
}

func (chain *Chain) String() string {
	bs := bytes.NewBufferString("")
	for _, block := range chain.blocks {
		bs.WriteString(block.String())
	}
	return bs.String()
}
