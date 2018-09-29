// +build ignore

package main

import(
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

func main(){
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil{
		panic(err)
	}

	request, err := http.NewRequest(
		"GET", "http://localhost:8888", nil)
	if err != nil{
		panic(err)
	}
	// リクエストを書き込んで送信
	request.Write(conn)
	// 読み取り
	response, err := http.ReadResponse(
		bufio.NewReader(conn), request)
		if err != nil{
			panic(err)
		}
		// 送られてきたレスポンスを読み込み
		dump, err := httputil.DumpResponse(response, true)
		if err != nil{
			panic(err)
		}
		// 表示
		fmt.Println(string(dump))
}