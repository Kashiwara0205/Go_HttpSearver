// +build ignore

package main
import(
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"io"
	"os"
	"compress/gzip"
)

func main(){
	// 送るメッセージは配列に格納しておく
	sendMessages := []string{
		"ASCII",
		"PROGRAMING",
		"PLUS",
	}
	current := 0
	var conn net.Conn = nil

	for {
		var err error
		// 初期設定
		if conn == nil{
			conn, err = net.Dial("tcp", "localhost:8888")
			if err != nil{
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}

		// POSTリクエスト生成
		request, err := http.NewRequest(
			"POST",
			"http://localhost:8888",
			strings.NewReader(sendMessages[current]))
		if err != nil{
			panic(err)
		}
		// gzipに対応していることを記述
		request.Header.Set("Accept-Encoding", "gzip")
		err = request.Write(conn)
		if err != nil{
			panic(err)
		}
		
		// レスポンス情報の入手
		response, err := http.ReadResponse(
			bufio.NewReader(conn), request)
		if err != nil{
			fmt.Println("Retry")
			// connをnilにして貼り直してもう一度接続しなおし
			conn = nil
			continue
		}
		// 情報が受け取れた場合出力
		dump, err := httputil.DumpResponse(response, true)
		if err != nil{
			panic(err)
		}
		fmt.Println(string(dump))

		// deferはこの関数が終わったら勝手に実行
		// rustのDROPと同じ
		defer response.Body.Close()
		// gzipで圧縮されているのであれば
		if response.Header.Get("Content-Encoding") == "gzip"{
			reader, err := gzip.NewReader(response.Body)
			if err != nil{
				panic(err)
			}
			// 圧縮されていた内容の出力
			io.Copy(os.Stdout, reader)
			reader.Close()
		}else{
			// 圧縮されていなくても出すので結局２回出すことになる
			io.Copy(os.Stdout, response.Body)
		}

		// 送信したメッセージ文だけカウント
		current++
		// 全部送信し終えたら終了
		if current == len(sendMessages){
			break
		}
	}
	conn.Close()
}