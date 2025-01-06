package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestResp(t *testing.T) {
	input := "+OK\r\n-Erradd\r\n*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n:55\r\n$5\r\nhello\r\n:99\r\n$3\r\nabc\r\n:-10\r\n"
	pr, pw := io.Pipe()
	go func() {
		for {
			for _, c := range input {
				pw.Write([]byte{byte(c)})
				time.Sleep(time.Millisecond * 100)
			}
			pw.Close()
		}
	}()

	reader := bufio.NewReader(pr)
	marshalBytes := make([]byte, 0)
	for {
		resp, err := Parse(reader)
		if err != nil {
			fmt.Println("Parse ERR:", err)
			break
		}
		bs, _ := resp.Marshal()
		marshalBytes = append(marshalBytes, bs...)
		fmt.Println("Parse RESP:", resp)
	}
	if !bytes.Equal([]byte(input), marshalBytes) {
		panic("marshal err")
	}
}
