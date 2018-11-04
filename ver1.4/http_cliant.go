// +build ignore

package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

func main(){
	sendMessages := []string{
		"ASCLL",
		"PROGRAMING",
		"PLUS",
	}
	current := 0
	var conn net.Conn = nil
	var err error
	requests := make(chan *http.Request, len(sendMessages))

	conn, err = net.Dial("tcp", "localhost:8888")
	if err != nil{
		panic(err)
	}
	fmt.Printf("Access: %d\n", current)

	// 連続で送信
	for i := 0; i < len(sendMessages); i++{
		lastMessage := i == len(sendMessages) - 1
		request, err := http.NewRequest(
			"GET",
			"http://localhost:8888?message" + sendMessages[i],
			nil)
		
		if lastMessage {
			request.Header.Add("Connection", "close")
		}else{
			request.Header.Add("Connection", "keep-alive")
		}
		if err != nil{
			panic(err)
		}
		// リクエスト送信
		err = request.Write(conn)
		if err != nil{
			panic(err)
		}

		fmt.Println("send: ", sendMessages[i])
		// requestをチャネルにつめる
		requests <- request
	}
	// いらないチャネルを閉じる
	close(requests)

	reader := bufio.NewReader(conn)
	// レスポンスを、ここでまとめて受信
	for request := range requests {
		// 受信
		response, err := http.ReadResponse(reader, request)
		if err != nil{
			panic(err)
		}
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		// サーバの情報を表示
		fmt.Println(string(dump))
		if current == len(sendMessages){
			break
		}
	}
}
