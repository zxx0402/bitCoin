package main

import (
	"bitCoin/someCompose"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	conn1, err1 := net.Dial("tcp", "127.0.0.1:20001")
	conn2, err2 := net.Dial("tcp", "127.0.0.1:20000")
	defer conn1.Close()
	defer conn2.Close()
	if err1 != nil || err2 != nil {
		log.Fatal(err1, err2)
	} else {
		chain := someCompose.NewBlockChain()
		marshal1, err1 := json.Marshal(chain)
		marshal2, err2 := json.Marshal(chain)
		if err1 != nil || err2 != nil {
			log.Fatal(err1, err2)
		} else {
			_, q1 := conn1.Write(marshal1)
			_, q2 := conn2.Write(marshal2)
			if q1 != nil || q2 != nil {
				log.Fatal("error", q1, q2)
				return
			}
		}
		var buf1 [1024]byte
		var buf2 [1024]byte
		n1, err5 := conn1.Read(buf1[:]) //将接收到的消息写到buf里面
		n2, err6 := conn2.Read(buf2[:]) //将接收到的消息写到buf里面
		if err5 != nil || err6 != nil {
			log.Fatal(err5, err6)
			return
		}
		fmt.Println("收到回复", string(buf1[:n1]))
		fmt.Println("收到回复", string(buf2[:n2]))
	}
}
