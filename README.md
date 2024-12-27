# resp
redis resp protocol

# test
```
func TestResp(t *testing.T) {
  input := "+OK\r\n-Erradd\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n:99\r\n$3\r\nabc\r\n:-10\r\n"
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
  if bytes.Compare([]byte(input), marshalBytes) != 0 { 
    panic("marshal err")
  }
}
```
# output
```
Parse RESP: {"Type":"+","Value":"OK"}
Parse RESP: {"Type":"-","Value":"Erradd"}
Parse RESP: {"Type":"*","Value":[{"Type":"$","Value":"foo"},{"Type":"$","Value":"bar"}]}
Parse RESP: {"Type":":","Value":99}
Parse RESP: {"Type":"$","Value":"abc"}
Parse RESP: {"Type":":","Value":-10}
Parse ERR: EOF
```
