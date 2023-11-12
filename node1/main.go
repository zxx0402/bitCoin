package main

import (
	"bitCoin/someCompose"
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"time"
)

var BlockChain someCompose.BlockChain
var victory string = "false"
var initVictory bool
var Trans []someCompose.Transaction
var mu sync.Mutex
var timeNow time.Time
var hasTransaction bool = false

func handleClient(clientNum int, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var buf [10000]byte
	for {
		n, err := reader.Read(buf[:])
		if err != nil {
		} else {
			s := string(buf[:n])
			if byte(s[2]) == 'B' {
				var bc someCompose.BlockChain
				err := json.Unmarshal(buf[:n], &bc)
				if err != nil {
					fmt.Printf("json trans err:%v\n", err)
					return
				}
				BlockChain = bc
				print("receive blockchain:::\n")
				someCompose.PrintBlockChain(&BlockChain)
				conn.Write([]byte("receive blockchain"))
				initVictory = true
			} else if byte(s[2]) == 'S' {
				hasTransaction = true
				var trans someCompose.Transaction
				json.Unmarshal(buf[:n], &trans) //接收到了信息，将信息添加到交易中
				Trans = append(Trans, trans)
				conn.Write([]byte("receive transaction")) //将接收到的数据回写
			}
			writer.Flush()
		}
	}
}
func Valid(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	println("i am waiting for block to valid")
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf = make([]byte, 1024)
		var result string
		for {
			n, err := reader.Read(buf)
			victory = "true"
			if err == io.EOF {
				break
			} else {
				result = result + string(buf[0:n])
				if n != cap(buf) {
					break
				}
			}
		}
		var timeAndBlock *someCompose.SendTimeAndBlock
		err := json.Unmarshal([]byte(result), &timeAndBlock)
		var block *someCompose.Block
		block = timeAndBlock.Block
		tt := timeAndBlock.Time
		if err == nil {
			fmt.Printf("获取到了块%v,time:%v\n", block.Nonce, time.Now())
			pow := someCompose.NewProofOfWork(block)
			valid := pow.IsValid(BlockChain) //验证
			someCompose.PrintBlock(block)
			if tt.Before(timeNow) == true {
				//表示对方挖到矿的时间早于我
				if valid == true {
					victory = "over"
					mu.Lock()
					BlockChain.Blocks = append(BlockChain.Blocks, block)
					mu.Unlock()
					println("验证成功：", block.Nonce)
					updateTransaction := someCompose.DeleteTransaction(Trans, block.Transaction)
					Trans = updateTransaction
					conn.Write([]byte("Validation victory!"))
				} else {
					victory = "false"
					println("验证失败: ", block.Nonce)
					conn.Write([]byte("Validation failed!"))
				}
				timeNow = someCompose.InitTime()
			} else {
				print("我先挖到矿")
				conn.Write([]byte("sorry i refused your mined!"))
				timeNow = someCompose.InitTime()
			}
		}
		victory = "false"
	}
}
func Mining(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	//挖矿
	println("i prepare mining")
	for {
		if initVictory == true && hasTransaction == true {
			lastOne := BlockChain.Blocks[len(BlockChain.Blocks)-1]
			block := someCompose.NewBlock1("newBlock", lastOne.Hash)
			pow := someCompose.NewProofOfWork(block)
			var nonce uint64
			var hash [32]byte
			var isMined bool = false
			for {
				if victory == "false" {
					fmt.Printf("挖矿中：%x\r", hash[:])
					data := pow.PrepareData(nonce)
					hash = sha256.Sum256(data)
					tempInt := new(big.Int)
					tempInt.SetBytes(hash[:])
					if tempInt.Cmp(pow.Target) == -1 {
						fmt.Printf("挖矿成功，hash：%x;nonce:%d\n", hash[:], nonce)
						isMined = true
						break
					} else {
						nonce++
					}
				} else if victory == "over" {
					break
				}
			}
			block.Hash = hash[:]
			block.Nonce = nonce
			result := someCompose.ValidNonce(BlockChain, block.Nonce)
			if victory == "false" && result == true && isMined == true {
				fmt.Printf("i am node2. i has mained,the nonce:%v\n,time:%v\n", nonce, time.Now())
				block.Transaction = Trans
				timeNow = time.Now()
				timeAndBlock := someCompose.NewSendTimeAndBlock(block, timeNow)
				marshal, _ := json.Marshal(timeAndBlock)
				_, err := conn.Write(marshal)
				if err != nil {
					fmt.Printf("err:%v\n", err)
				} else {
					fmt.Printf("将挖矿结果发送给另一个节点：%v\n", block.Nonce)
					var buf [128]byte
					n, _ := conn.Read(buf[:])
					if string(buf[:n]) == "Validation victory!" {
						mu.Lock()
						BlockChain.Blocks = append(BlockChain.Blocks, block)
						mu.Unlock()
						updateTransaction := someCompose.DeleteTransaction(Trans, block.Transaction)
						Trans = updateTransaction
						fmt.Printf("另一个节点检验成功，nonce：%v\n", block.Nonce)
						print("validation victory and add to my chain\n")
						fmt.Printf("len(chain):%v\n", len(BlockChain.Blocks))
						someCompose.PrintBlockChain(&BlockChain)
					} else if string(buf[:n]) == "Validation failed!" {
						fmt.Printf("挖矿结果未经过另一个节点的检验，nonce：%v\n", block.Nonce)
						fmt.Println("Validation false!")
						continue
					}
				}
				timeNow = someCompose.InitTime()
			}
		}
	}
}
func main() {
	var wg sync.WaitGroup
	timeNow = someCompose.InitTime()
	serverAddress := "127.0.0.1:20000"
	serverAddress1 := "127.0.0.1:20003"
	listener, err := net.Listen("tcp", serverAddress)
	listen1, err := net.Listen("tcp", serverAddress1) //指明使用的网络和地址
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()
	defer listen1.Close()
	var conn net.Conn
	var err1 error
	for {
		conn, err1 = net.Dial("tcp", "127.0.0.1:20002")
		if err1 != nil {
			//fmt.Printf("conn:%v\n", err1)
		} else {
			wg.Add(1)
			go Mining(conn, &wg)
			break
		}
	}
	fmt.Println("Server is listening on", serverAddress)
	accept1, _ := listen1.Accept()
	wg.Add(1)
	go Valid(accept1, &wg)
	var clientNum int
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		clientNum++
		fmt.Printf("Client %d connected\n", clientNum)
		wg.Add(1)
		go handleClient(clientNum, conn, &wg)
	}
	wg.Wait()
}
