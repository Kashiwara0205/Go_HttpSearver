// +build ignore

package main
import(
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func main(){
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil{
		panic(err)
	}
	defer conn.Close()
	request, err := http.NewRequest(
		"GET",
		"http://localhost:8888",
		nil)
	if err != nil{
		panic(err)
	}
	err = request.Write(conn)
	if err != nil{
		panic(err)
	}
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	if err != nil{
		panic(err)
	}
	dump, err := httputil.DumpResponse(response, false)
	if err != nil{
		panic(err)
	}

	// サーバ側からのメッセージを表示
	fmt.Println(string(dump))

	// TransferEncodingがチャンクじゃなかったらエラー
	// チャンクが無くても動くけどTransferEncodingはチャンクじゃないとダメ
	if len(response.TransferEncoding) < 1 ||
		response.TransferEncoding[0] != "chunked"{
		panic("wrong transfer encoding")
	}

	// 長い、あの文章はここでいイテレートされて表示される
	for {
		// \nまでのサイズ取得
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF{
			break
		}
		// 16進数の64bitからInt型の数値に置き換えてサイズを取得
		size, err := strconv.ParseInt(
			string(sizeStr[:len(sizeStr)-2]), 16, 64)
		// サイズを表す数値で0が送られてきたら終了
		if size == 0{
			break
		}
		if err != nil{
			panic(err)
		}
		// byte型配列を送られてきたサイズの分だけ生成
		line := make([]byte, int(size))
		// size分だけ読み込み
		io.ReadFull(reader, line)
		reader.Discard(2)
		// 表示
		fmt.Printf("%d bytes: %s\n", size, string(line))
	}
}