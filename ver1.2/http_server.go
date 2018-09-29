package main

import(
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// gzipを受け入れる事が出来るクライアントなのかどうかを判断
func isGzipAcceptable(request *http.Request) bool{
	return strings.Index(
		strings.Join(request.Header["Accept-Encoding"], ","),
		"gzip") != -1
}

func processSession(conn net.Conn){
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	defer conn.Close()
	for {
		// 今から５秒後にレスポンスが帰ってこなかったら切断
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// request取得
		request, err := http.ReadRequest(bufio.NewReader(conn))

		if err != nil{
			neterr, ok := err.(net.Error)
			// 5秒まったけど帰ってこなかった
			if ok && neterr.Timeout(){
				fmt.Println("Timeout")
				break
			}else if err == io.EOF{
				// メッセージの終わり
				break
			}
		}
		// requestを取得
		dump, err := httputil.DumpRequest(request, true)
		// rustの.excpetみたいにかけたら楽なんだが
		if err != nil{
			panic(err)
		}
		// 送られてきた内容を表示
		fmt.Println(string(dump))

		// サーバ側からレスポンスを返却
		response := http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: make(http.Header),
		}

		// gzipが受け入れられるクライアントかどうか検証
		if isGzipAcceptable(request){
			content := "Hello World (gzipped)\n"
			var buffer bytes.Buffer
			// 送る内容をgzipで圧縮
			writer := gzip.NewWriter(&buffer)
			io.WriteString(writer, content)
			writer.Close()
			response.Body = ioutil.NopCloser(&buffer)
			response.ContentLength = int64(buffer.Len())
			response.Header.Set("Content-Encoding", "gzip")
		}else{
			// 対応してなかったら、そのまま送信
			content := "Hello world\n"
			response.Body = ioutil.NopCloser(
				strings.NewReader(content))
			response.ContentLength = int64(len(content))
		}
		response.Write(conn)
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