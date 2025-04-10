package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine 
}

type RequestLine struct {
	HttpVersion string 
	RequestTarget string 
	Method string 
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	res,err := parseRequestLine(string(b))
	if err != nil {
		return nil,err 
	}


	

	return &Request{
		*res,
	},nil 
}

func parseRequestLine(data string) (*RequestLine,error){
	
	lines := strings.Split(data,"\r\n")
	
	requestLine := strings.Split(lines[0]," ")
	if len(requestLine) != 3 {
		return nil,fmt.Errorf("invalid request line: %s",requestLine)
	}
	
	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(requestLine[0]) {
		return nil,fmt.Errorf("invalid method: %s",requestLine[0])
		
	}
	if strings.Compare(requestLine[2],"HTTP/1.1") != 0 {
		return nil,fmt.Errorf("invalid http version: %s",requestLine[2])
	}
	method := requestLine[0]
	targetPath := requestLine[1]
	httpVersion := "1.1"
	
	return &RequestLine{
		HttpVersion: httpVersion,
		RequestTarget: targetPath,
		Method: method,},nil 
}