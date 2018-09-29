package main

import(
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

// 送信する内容を文字列型配列で用意
var contents = []string{
	"これは、わたしが小さいときに、村の茂兵（もへい）というおじいさんからきいたお話です。",
	"むかしは、わたしたちの村のちかくの、中山というところに小さなお城（しろ）があって、",
	"中山さまというおとのさまがおられたそうです。",
	"その中山から、すこしはなれた山の中に、「ごんぎつね」というきつねがいました。",
	"ごんは、ひとりぼっちの小ぎつねで、しだのいっぱいしげった森の中に穴（あな）をほって住んでいました。",
}

func processSession(conn net.Conn){
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	defer conn.Close()
	for{
		// リクエスト待ち
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil{
			// 読み取ってきた内容が終わりなら終了
			if err == io.EOF{
				break
			}
			panic(err)
		}
		dump, err := httputil.DumpRequest(request, true)
		if err != nil{
			panic(err)
		}
		// リクエスト内容を表示
		fmt.Println(string(dump))

		// 送信
		fmt.Fprintf(conn, strings.Join([]string{
			"HTTP/1.1 200 OK",
			"Content-Type: text/plain",
			"Transfer-Encoding: chunked", // チャンク形式で送ることを意味する
			"","",
		}, "\r\n"))

		// 一つ一つbyte形式にして送る
		for _, content := range contents{
			bytes := []byte(content)
			// %xは16進　故に16進数でサイズを管理していることとなる
			fmt.Fprintf(conn, "%x\r\n%s\r\n", len(bytes), content) // バイトの数はヘッダから外れ、こちらで別途送信
		}
		// 最後に0を数値を送信
		fmt.Fprintf(conn, "0\r\n\r\n")
	}
}

func main(){
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil{
		panic(err)
	}
	fmt.Println("Server is running at localhost:8888")
	for{
		conn, err := listener.Accept()
		if err != nil{
			panic(err)
		}
		go processSession(conn)
	}
}