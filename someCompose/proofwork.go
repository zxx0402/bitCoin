package someCompose

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//实现挖矿功能

type ProofOfWork struct {
	Block  *Block
	Target *big.Int //这个结构提供了很多方法：比较、把哈希值设置为big.Int类型
}

// 创建ProofOfWork
// block由用户提供
// target目标值由系统提供
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		Block: block,
	}
	//难度值先写死
	targetstr := "0000100000000000000000000000000000000000000000000000000000000000" //64位
	tempBigInt := new(big.Int)
	//将难度值赋值给bigInt
	tempBigInt.SetString(targetstr, 16)
	pow.Target = tempBigInt
	return &pow
}

// 计算哈希值，不断变化nonce，使得sha256（数据）+nonce<难度值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	fmt.Println("开始挖矿。。。")
	for {
		fmt.Printf("%x\r", hash[:])
		//1、拼接字符串
		data := pow.PrepareData(nonce)
		//2、sha256
		hash = sha256.Sum256(data)
		//将hash转换为bigint类型
		tempInt := new(big.Int)
		tempInt.SetBytes(hash[:])
		//3、比较当前的哈希值与难度值
		//if 当前哈希<难度值 return 哈希，nonce
		//else nonce++
		//当前的计算的哈希值.Cmp(难度值)
		if tempInt.Cmp(pow.Target) == -1 {
			fmt.Printf("挖矿成功，hash：%x\n;nonce:%d\n", hash[:], nonce)
			break
		} else {
			nonce++
		}
	}
	return hash[:], nonce
}
func (pow *ProofOfWork) Run1(valid string) ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	fmt.Println("开始挖矿。。。")
	fmt.Println("valid", valid)
	for {
		if valid == "false" {
			//表示节点未接收到其他区块，所以持续挖矿
			fmt.Printf("%x\r", hash[:])
			//1、拼接字符串
			data := pow.PrepareData(nonce)
			//2、sha256
			hash = sha256.Sum256(data)
			//将hash转换为bigint类型
			tempInt := new(big.Int)
			tempInt.SetBytes(hash[:])
			//3、比较当前的哈希值与难度值
			//if 当前哈希<难度值 return 哈希，nonce
			//else nonce++
			//当前的计算的哈希值.Cmp(难度值)
			if tempInt.Cmp(pow.Target) == -1 {
				fmt.Printf("挖矿成功，hash：%x\n;nonce:%d\n", hash[:], nonce)
				break
			} else {
				nonce++
			}
		} else if valid == "over" {
			//表示其他节点已经挖到矿了，所以要结束挖矿，重新开始
			println("Mining over")
			break
		} else {
			//表示接收到了其他节点发来的消息，所以暂时停止挖矿
			//println("i am waiting")
		}
	}
	return hash[:], nonce
}

// 拼接nonce和block数据
func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {
	b := pow.Block
	tmp := [][]byte{
		uintToByte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(nonce),
		//b.Hash,
		b.Data,
	}
	data := bytes.Join(tmp, []byte{})
	return data
}

// 其他区块验证
func (pow *ProofOfWork) IsValid(bc BlockChain) bool {
	//1、获取区块
	//2、拼接数据
	data := pow.PrepareData(pow.Block.Nonce)
	//3、计算sha256
	hash := sha256.Sum256(data)
	//4、与难度值比较
	tempInt := new(big.Int)
	tempInt.SetBytes(hash[:])
	lastOne := len(bc.Blocks) - 1
	theLastBlockHash := bc.Blocks[lastOne].Hash
	h := pow.Block.PrevHash
	if tempInt.Cmp(pow.Target) == -1 && string(theLastBlockHash) == string(h) {
		return true
	} else {
		return false
	}
}
func (pow *ProofOfWork) IsValid1() bool {
	//1、获取区块
	//2、拼接数据
	data := pow.PrepareData(pow.Block.Nonce)
	//3、计算sha256
	hash := sha256.Sum256(data)
	//4、与难度值比较
	tempInt := new(big.Int)
	tempInt.SetBytes(hash[:])
	if tempInt.Cmp(pow.Target) == -1 {
		return true
	} else {
		return false
	}
}
