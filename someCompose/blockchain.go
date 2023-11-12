package someCompose

import "fmt"

// 定义区块链结构(使用数组模拟)
type BlockChain struct {
	Blocks []*Block //区块链
}

// 创世语
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// 创建区块链 创建blockchain，同时添加一个创世快
func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock(genesisInfo, nil)
	bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}
	return &bc
}

// 添加区块
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	//创建block
	newBlock := NewBlock(data, lastBlock.Hash)
	//添加到bc中
	bc.Blocks = append(bc.Blocks, newBlock)

}
func PrintBlockChain(bc *BlockChain) {
	for i, block := range bc.Blocks {
		fmt.Printf("++++++当前区块高度++++：%d\n", i)
		//fmt.Printf("Version：%d\n", block.Version)
		fmt.Printf("PreHash：%x\n", block.PrevHash)
		//fmt.Printf("MerkelRoot：%x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp：%d\n", block.TimeStamp)
		//fmt.Printf("Bits：%d\n", block.Bits)
		fmt.Printf("Nonce：%d\n", block.Nonce)
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Trans:%d\n", len(block.Transaction))
		//fmt.Printf("Hash：%x\n", block.Hash)
		//fmt.Printf("Data：%s\n", block.Data)
		//pow := NewProofOfWork(block)
		//fmt.Printf("IsValid：%v\n", pow.IsValid())
	}
}
func ValidNonce(bc BlockChain, nonce uint64) (result bool) {
	for _, block := range bc.Blocks {
		if block.Nonce == nonce {
			return false
		}
	}
	return true
}
