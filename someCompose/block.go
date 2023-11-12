package someCompose

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	//版本号
	Version uint64
	//前区块哈希
	PrevHash []byte
	//交易的根哈希
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值--系统提供一个数值，用于计算出哈希值
	Bits uint64
	//随机数，挖矿要求的数值
	Nonce uint64
	//哈希
	Hash []byte
	//数据
	Data []byte
	//交易
	Transaction []Transaction
}

func NewBlock(data string, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		MerkleRoot: nil,
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0, //随意写的
		Nonce:      0, //随意写的
		PrevHash:   prevHash,
		Hash:       nil,
		Data:       []byte(data),
	}
	//计算哈希
	//b.setHash()
	pow := NewProofOfWork(&b)
	//挖矿
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce
	return &b
}
func NewBlock1(data string, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		MerkleRoot: nil,
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0, //随意写的
		Nonce:      0, //随意写的
		PrevHash:   prevHash,
		Hash:       nil,
		Data:       []byte(data),
	}
	//计算哈希
	//b.setHash()
	//pow := NewProofOfWork(&b)
	////挖矿
	//hash, nonce := pow.Run()
	//b.Hash = hash
	//b.Nonce = nonce
	return &b
}

// 计算区块哈希值
func (b *Block) setHash() {
	//sha256
	//data是block各个字段拼成的字节流
	tmp := [][]byte{
		uintToByte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(b.Nonce),
		b.Hash,
		b.Data,
	}
	data := bytes.Join(tmp, []byte{}) //将三个切片拼接在一起
	hash := sha256.Sum256(data)
	b.Hash = hash[:]

}
func PrintBlock(block *Block) {
	println("--------------------------------------------------")
	fmt.Printf("PreHash：%x\n", block.PrevHash)
	//fmt.Printf("MerkelRoot：%x\n", block.MerkleRoot)
	//fmt.Printf("TimeStamp：%d\n", block.TimeStamp)
	//fmt.Printf("Bits：%d\n", block.Bits)
	fmt.Printf("Nonce：%d\n", block.Nonce)
	fmt.Printf("Hash：%x\n", block.Hash)
	fmt.Printf("Trans:%d\n", len(block.Transaction))
	println("-------------------------------------------------------")
	//fmt.Printf("Hash：%x\n", block.Hash)
	//fmt.Printf("Data：%s\n", block.Data)
	//pow := NewProofOfWork(block)
	//fmt.Printf("IsValid：%v\n", pow.IsValid())
}
