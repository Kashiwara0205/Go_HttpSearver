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
 
func writeToConn(sessionResponses chan chan *http.Response, conn net.Conn){
	defer conn.Close()
	// 順番にsessionを取り出す
	for sessionResponses := range sessionResponses{
		response := <- sessionResponses
		response.Write(conn)
		// 	不必要になったチャネルを閉じる
		close(sessionResponses)
	}
}

func handleRequest(request *http.Request,
				   resultReceiver chan *http.Response){
	dump, err := httputil.DumpRequest(request, true)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(dump))
	content := "Hello World\n"

	response := &http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
		ContentLength: int64(len(content)),
		Body: ioutil.NopCloser(strings.NewReader(content)),
	}

	// こいつが走ったら上のwriteToConnメソッドのresponse.Write(conn)が走りだす
	resultReceiver <- response
}

func processSession(conn net.Conn){
	fmt.Printf("Acept %v \n", conn.RemoteAddr())
	// chan chanは同期的な処理をするためのもの
	// AのチャネルBチャネルにアクセスしにいくみたいな感じ
	sessionResponses := make(chan chan *http.Response, 50)
	defer close(sessionResponses)
	// sessionResponsesが同期処理の受け渡し窓口
	go writeToConn(sessionResponses, conn)

	reader := bufio.NewReader(conn)
	for {
		// ココらへんはKeep-Aliveの仕組みと同じである
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		// request受け取り
		request, err := http.ReadRequest(reader)
		if err != nil{
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout(){
				fmt.Println("Timeout")
				break
			}else if err == io.EOF{
				break
			}
			panic(err)
		}
		// 新たにチャネル作成
		sessionResponse := make(chan *http.Response)
		// チャネルのチャネル　sessionResponseからsessionResponsesが送られてくる用に設定
		sessionResponses <- sessionResponse

		// requestを非同期にさばいていく
		go handleRequest(request, sessionResponse)
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