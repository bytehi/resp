package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

func (r *RESP) Marshal() ([]byte, error) {
	bs, err := r.XXX_Marshal(nil)
	if err != nil {
		return nil, err
	}
	return bs.Bytes(), nil
}

func (r *RESP) XXX_Marshal(buf *bytes.Buffer) (*bytes.Buffer, error) {
	if buf == nil {
		size, err := r.XXX_Size()
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(make([]byte, 0, size))
	}

	writeByte := buf.WriteByte
	writeBytes := buf.Write
	writeInt := func(val int64) {
		s := strconv.FormatInt(val, 10)
		writeBytes([]byte(s))
	}

	switch r.Type {
	case SimpleString, Error:
		writeByte(byte(r.Type))
		writeBytes([]byte(r.Value.(string)))
		writeByte('\r')
		writeByte('\n')
	case BulkString:
		bs := []byte(r.Value.(string))
		writeByte(byte(r.Type))
		writeInt(int64(len(bs)))
		writeByte('\r')
		writeByte('\n')
		writeBytes(bs)
		writeByte('\r')
		writeByte('\n')
	case Integer:
		number := r.Value.(int64)
		writeByte(byte(r.Type))
		writeInt(number)
		writeByte('\r')
		writeByte('\n')
	case Array:
		members := r.Value.([]*RESP)
		writeByte(byte(r.Type))
		writeInt(int64(len(members)))
		writeByte('\r')
		writeByte('\n')
		for _, member := range members {
			if _, err := member.XXX_Marshal(buf); err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unknown type:%s value:%v", r.Type, r.Value)
	}
	return buf, nil
}

func intToByteSize(val int64) int {
	str := strconv.FormatInt(val, 10)
	return len(str)
}

func (r *RESP) XXX_Size() (int, error) {
	switch r.Type {
	case SimpleString, Error:
		strLen := len([]byte(r.Value.(string)))
		return 1 + strLen + 2, nil
	case BulkString:
		strLen := len([]byte(r.Value.(string)))
		return 1 + intToByteSize(int64(strLen)) + 2 + strLen + 2, nil
	case Integer:
		number := r.Value.(int64)
		return 1 + intToByteSize(number) + 2, nil
	case Array:
		members := r.Value.([]*RESP)
		size := 1 + intToByteSize(int64(len(members))) + 2
		for _, member := range members {
			size1, err := member.XXX_Size()
			if err != nil {
				return 0, err
			}
			size += size1
		}
		return size, nil
	default:
		return 0, fmt.Errorf("unknown type:%s value:%v", r.Type, r.Value)
	}
}
