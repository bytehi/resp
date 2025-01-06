package resp

import (
	"encoding/json"
	"fmt"
)

type RespType byte

func (r RespType) String() string {
	return fmt.Sprintf("\"%c\"", r)
}

func (r RespType) MarshalJSON() ([]byte, error) {
	return []byte(r.String()), nil
}

const (
	SimpleString RespType = '+' // Simple string response, value is string
	Error        RespType = '-' // Error response, value is string
	Integer      RespType = ':' // Integer response, value is int64
	BulkString   RespType = '$' // Bulk string response, value is string
	Array        RespType = '*' // Array response, value is []*RESP
)

type RESP struct {
	Type  RespType    `json:"type"`
	Value interface{} `json:"value"`
}

func (r *RESP) String() string {
	datas, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("error marshaling response: %v", err)
	}
	return string(datas)
}

// create a simple string response
func NewSimpleString(value string) *RESP {
	return &RESP{Type: SimpleString, Value: value}
}

// create a error response
func NewError(value string) *RESP {
	return &RESP{Type: Error, Value: value}
}

// create a integer response
func NewInteger(value int64) *RESP {
	return &RESP{Type: Integer, Value: value}
}

// create a bulk string response
func NewBulkString(value string) *RESP {
	return &RESP{Type: BulkString, Value: value}
}

// create a array response
func NewArray(value []*RESP) *RESP {
	return &RESP{Type: Array, Value: value}
}
