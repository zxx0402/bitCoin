package main

import (
	"bitCoin/someCompose"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

/*
客户端，每隔10ms向两个node，即两个服务端发送交易信息
*/
func init() {
	rand.Seed(time.Now().UnixMicro()) //设置种子
}

func main() {
	tick := time.Tick(time.Millisecond * 10)
	//建立连接
	conn1, err1 := net.Dial("tcp", "127.0.0.1:20001")
	conn2, err2 := net.Dial("tcp", "127.0.0.1:20000")
	defer conn1.Close()
	defer conn2.Close()
	for _ = range tick {
		if err1 != nil || err2 != nil {
		} else {
			//为交易值设置随机数
			intn := rand.Intn(10)
			//创建新的交易
			trans := someCompose.NewTransaction("client", "receiver", intn)
			trans1 := someCompose.NewTransaction("client", "receiver", intn)
			marshal1, err3 := json.Marshal(trans)
			marshal2, err4 := json.Marshal(trans1)
			if err3 != nil || err4 != nil {
				fmt.Printf("err3:%v\n", err3)
				fmt.Printf("err4:%v\n", err4)
			} else {
				_, q1 := conn1.Write(marshal1)
				_, q2 := conn2.Write(marshal2)
				if q1 != nil || q2 != nil {
					fmt.Printf("q1:%v\n", q1)
					fmt.Printf("q2:%v\n", q2)
					return
				}
				var buf1 [1024]byte
				var buf2 [1024]byte
				n1, err5 := conn1.Read(buf1[:]) //将接收到的消息写到buf里面
				n2, err6 := conn2.Read(buf2[:]) //将接收到的消息写到buf里面
				println("--------------------------")
				if err5 != nil || err6 != nil {
					fmt.Printf("err5:%v\n", err5)
					fmt.Printf("err6:%v\n", err6)
					return
				}
				fmt.Println("收到回复", string(buf1[:n1]))
				fmt.Println("收到回复", string(buf2[:n2]))

			}
		}
	}
}
