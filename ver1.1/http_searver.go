package main

import(
	"bufio"
	"io"
	"time"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main(){
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

		// Keep-Alive対応版
		go func(){
			defer conn.Close()
			for{
				// タイムアウトの設定(5秒待ち)
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				// リクエストの読み込み
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil{
					fmt.Println(err)
					// タイムアウト、クローズ時に終了
					// それ以外はエラー
					neterr, ok := err.(net.Error)
					if ok && neterr.Timeout(){
						fmt.Println("Timeout")
						break
					}else if err == io.EOF{
						break
					}
					panic(err)
				}

				// リクエストの表示
				dump, err := httputil.DumpRequest(request, true)
				if err != nil{
					panic(err)
				}
				// http1.0からはチャンクにcontent-lengthが含まれるようになったぞ！
				fmt.Println(string(dump))
				content := "Hello World\n"

				// レスポンスの書き込み
				response := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 1,
					ContentLength: int64(len(content)), // バイト数が付与されたぞ！
					Body: ioutil.NopCloser(
						  strings.NewReader(content)),
				}
				response.Write(conn)
			}
		}()
	}
}