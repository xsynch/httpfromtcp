package headers

import (
	"bytes"
	"fmt"
	
	"regexp"
	"strings"
)


type Headers map[string]string

type State int 
const (
	initialized State = iota
	done 
)
const (
bufferSize = 8
)


type RequestHeaders struct {
	HeaderLine Headers
	HTTPReadStatus State  
}



type HeaderLine struct {
	FieldName string 
	FieldValue string 
}



// func RequestFromReader(reader io.Reader) (*Request, error) {
// 	readToindex := 0
// 	buf := make([]byte,bufferSize)

// 	req := &RequestHeaders{
// 		HTTPReadStatus: initialized,
// 	}
	
	
// 	for req.HTTPReadStatus != done{
		
// 		if readToindex == len(buf){
// 			newBuf := make([]byte,len(buf) * 2)
// 			_ = copy(newBuf,buf)
// 			buf = newBuf

// 		}
// 		b,err := reader.Read(buf[readToindex:])
// 		if err != nil && err != io.EOF{
// 			return nil, err				
// 		}
// 		readToindex += b
		
// 		bytesConsumed,status,err := req.HeaderLine.Parse((buf[:readToindex]))
// 		if err != nil {
// 			return nil, err 
// 		}
// 		if bytesConsumed > 0 {
// 			copy(buf,buf[bytesConsumed:readToindex])
// 			readToindex -= bytesConsumed
// 		}

// 		if err == io.EOF{
// 			req.HTTPReadStatus = done
// 			break
// 		}
				
// 	}
	

// 	return req,nil 
// }


// func parseHeaderLine(data []byte) (Headers,int ,error){
// 	if !bytes.Contains(data,[]byte("\r\n")) {
// 		return nil, 0,nil 
// 	}
// 	bytesConsumed := bytes.Index(data, []byte("\r\n")) + 2
	
// 	lines := strings.Split(string(data),"\r\n")
	
// 	headerLine := strings.Split(lines[0],": ")
// 	if len(headerLine) != 3 {
// 		return nil, 0,fmt.Errorf("invalid request line: %s",headerLine)
// 	}
	
// 	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(headerLine[0]) {
// 		return nil, 0,fmt.Errorf("invalid method: %s",headerLine[0])
		
// 	}
// 	httpVer := strings.Split(headerLine[2],"/")
// 	if len(httpVer) !=2 || httpVer[0] != "HTTP" || httpVer[1] != "1.1"{
// 		return nil, 0,fmt.Errorf("invalid http version: %s",headerLine[2])
// 	}
	

	
// 	return Headers{},bytesConsumed,nil 
// }

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	if len(data) == 0 {
		return 0,false,nil 
	}

	if bytes.HasPrefix(data,[]byte("\r\n")){
		return 2, true,nil 
	}

	idx := bytes.Index(data,[]byte("\r\n"))
	if idx == -1 {
		return 0,false,nil
	}

	line := data[:idx]

	colonIdx := bytes.Index(line,[]byte(":"))
	if colonIdx == -1 {
		return 2,false,fmt.Errorf("invalid header, no colon found in %s",string(line))
	}
	key := string(line[:colonIdx])
	if string(bytes.TrimRight(line[:colonIdx]," ")) != key {
		return 0, false,fmt.Errorf("invalid header, spaces before colon: %s",string(line[:colonIdx]))
	}

	key = strings.TrimSpace(key)
	val := string(bytes.TrimSpace((line[colonIdx+1:])))
	if len(val) < 1 {
		return 0,false,fmt.Errorf("header value must contain a value: %s",val)
	}
	if isValidHeader(key) {
		key = strings.ToLower(key)
		v,ok := h[key]
		if ok {			
			h[key] = fmt.Sprintf("%s, %s",v, val)
			
		} else {
			h[key] = val
		}
	} else {
		return 0,false,fmt.Errorf("key contains invalid character: %s",key)
	}
	return idx + 2,false,nil 



	

}

func isValidHeader(word string) bool {
	pattern := "^[a-zA-Z0-9\\!#\\$\\%&'\\*\\+-\\.^_`|~]*$"
	// escapedPattern := regexp.QuoteMeta(pattern)
	regex, _ := regexp.Compile(pattern)
	return regex.MatchString(word)
}

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Get(key string) (string,error) {
	val,ok := h[strings.ToLower(key)]
	if ok {
		return val,nil
	} else {
		return "",fmt.Errorf("header not found: %s",val)
	}
	

}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}