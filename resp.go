package resp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type RESP_TYPE byte

func (r RESP_TYPE) String() string {
	return fmt.Sprintf("\"%c\"", r)
}

func (r RESP_TYPE) MarshalJSON() ([]byte, error) {
	return []byte(r.String()), nil
}

const (
	SimpleString RESP_TYPE = '+' //Value is string
	Error        RESP_TYPE = '-' //Value is string
	Integer      RESP_TYPE = ':' //value is int64
	BulkString   RESP_TYPE = '$' //value is string
	Array        RESP_TYPE = '*' //value is []*RESP
)

type RESP struct {
	Type  RESP_TYPE
	Value interface{}
}

func (r *RESP) String() string {
	datas, _ := json.Marshal(r)
	return string(datas)
}

func Parse(reader *bufio.Reader) (*RESP, error) {
	readline := func() ([]byte, error) {
		line, isPrefix, err := reader.ReadLine()
		for isPrefix && err == nil {
			var line2 []byte
			line2, isPrefix, err = reader.ReadLine()
			if err == nil {
				line = append(line, line2...)
			}
		}
		if err != nil {
			return nil, err
		}
		return line, nil
	}

	firstLine, err := readline()
	if err != nil {
		return nil, err
	}
	typ := RESP_TYPE(firstLine[0])
	switch typ {
	case SimpleString, Error:
		return &RESP{Type: typ, Value: string(firstLine[1:])}, nil
	case Integer:
		val, err := strconv.ParseInt(string(firstLine[1:]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &RESP{Type: typ, Value: val}, nil
	case BulkString:
		strLen, err := strconv.Atoi(string(firstLine[1:]))
		if err != nil {
			return nil, err
		}
		buf := make([]byte, strLen+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}
		return &RESP{Type: typ, Value: string(buf[:strLen])}, nil
	case Array:
		num, err := strconv.Atoi(string(firstLine[1:]))
		if err != nil {
			return nil, err
		}
		members := make([]*RESP, num)
		for i := 0; i < num; i++ {
			resp2, err := Parse(reader)
			if err != nil {
				return nil, err
			}
			members[i] = resp2
		}
		return &RESP{
			Type:  Array,
			Value: members,
		}, nil
	default:
		return nil, fmt.Errorf("unknown type:%s", typ)
	}
}
