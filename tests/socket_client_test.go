package tests

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"jim/client/core"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	HOST string = "localhost"
	PORT int    = 4009
	conn net.Conn
)

func TestSocketClient(t *testing.T) {
	var err error
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		return
	}
	go loopRead()
	//loopWrite()
	//inputer()
	str:="hello world"
	data, err:=core.Encode([]byte(str))
	if err!=nil{
		return
	}
	conn.Write(data)
	time.Sleep(3*time.Second)
}

func loopRead() {
	for {
		inputReader := bufio.NewReader(conn)
		data, err := core.Decode(inputReader)
		if err != nil {
			if io.EOF == err {
				panic(errors.New("connection is closed"))
			}
			continue
		}
		fmt.Println(string(data))
	}
}

func loopWrite() {
	input := bufio.NewReader(os.Stdin)
	for {
		str, err := input.ReadString('\n')
		if err != nil {
			continue
		}
		str = str[0 : len(str)-1]
		str = strings.Trim(str, " ")
		data, err := core.Encode([]byte(str))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		conn.Write(data)
	}
}

func inputer() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Printf("Please type in something:\n")
	for input.Scan() {
		line := input.Text()
		if line == "\n" {
			fmt.Println(line)
			break
		}
	}
}