package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

// Parse reads a RESP message from reader and returns the parsed RESP value
func Parse(reader *bufio.Reader) (*RESP, error) {
	// Read the type byte
	typ, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typ {
	case byte(SimpleString), byte(Error):
		str, err := readLine(reader, 0)
		if err != nil {
			return nil, err
		}
		return &RESP{Type: RespType(typ), Value: str}, nil

	case byte(Integer):
		str, err := readLine(reader, 20)
		if err != nil {
			return nil, err
		}
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		return &RESP{Type: Integer, Value: i}, nil

	case byte(BulkString):
		str, err := readLine(reader, 20)
		if err != nil {
			return nil, err
		}
		length, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		if length == -1 {
			return &RESP{Type: BulkString, Value: ""}, nil
		}
		buf := make([]byte, length+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}
		return &RESP{Type: BulkString, Value: string(buf[:length])}, nil

	case byte(Array):
		str, err := readLine(reader, 0)
		if err != nil {
			return nil, err
		}
		length, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		if length == -1 {
			return &RESP{Type: Array, Value: []*RESP{}}, nil
		}
		arr := make([]*RESP, length)
		for i := range arr {
			resp, err := Parse(reader)
			if err != nil {
				return nil, err
			}
			arr[i] = resp
		}
		return &RESP{Type: Array, Value: arr}, nil

	default:
		return nil, fmt.Errorf("unknown RESP type: %c", typ)
	}
}

// readLine reads until CRLF and returns the line without CRLF
func readLine(reader *bufio.Reader, limit int) (string, error) {
	var line []byte
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		line = append(line, b...)
		if len(b) >= 2 && b[len(b)-2] == '\r' {
			line = line[:len(line)-2]
			return *(*string)(unsafe.Pointer(&line)), nil
		}
		if limit > 0 && len(line) >= limit {
			return "", fmt.Errorf("line too long: %d", len(line))
		}
	}
}
