package main

import(
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main(){
	// 8888番でポート待ち
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil{
		panic(err)
	}
	fmt.Println("Server is runnnig at localhost:8888")

	for {

		conn, err := listener.Accept()
		if err != nil{
			panic(err)
		}

		// 並行処理(ゴルーチン)
		go func(){
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			request, err := http.ReadRequest(
				bufio.NewReader(conn))
			if err != nil{
				panic(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil{
				panic(err)
			}
			fmt.Println(string(dump))
			// ここでクライアントサイドに出す時のやつが整形されてる
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body: ioutil.NopCloser(
					strings.NewReader("Hello World\n")),
			}
			// 返却
			response.Write(conn)
			// ver1.0のhttpは１回１回プツプツ切れてたらしい
			conn.Close()
		}()
	}
}