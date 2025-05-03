package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/xsynch/httpfromtcp/internal/headers"
)

type State int 
const (
	requestStateInitialized State = iota
	requestStateParsingHeaders 
	requestStateParsingBody
	requestStateDone 
)
const (
bufferSize = 8
)


type Request struct {
	RequestLine RequestLine
	HTTPReadStatus State 
	Headers headers.Headers
	Body []byte
	bodyLengthRead int 
	
}

type RequestLine struct {
	HttpVersion string 
	RequestTarget string 
	Method string 
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToindex := 0
	buf := make([]byte,bufferSize)

	req := &Request{
		HTTPReadStatus: requestStateInitialized,
		Headers: headers.NewHeaders(),
		Body: make([]byte, 0),
	}
	for req.HTTPReadStatus != requestStateDone{
		if readToindex == len(buf){
			newBuf := make([]byte,len(buf) * 2)
			_ = copy(newBuf,buf)
			buf = newBuf

		}
		b,err := reader.Read(buf[readToindex:])
		if err != nil {
			if errors.Is(err,io.EOF){
				if (req.HTTPReadStatus != requestStateDone){
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d",req.HTTPReadStatus,b)
				}
				break
			}
			return nil, err 						
		}

		readToindex += b
		
		bytesConsumed,err := req.parse(buf[:readToindex])
		if err != nil {
			return nil, err 
		}
		
		copy(buf,buf[bytesConsumed:readToindex])
		readToindex -= bytesConsumed
		
	}
	
	return req,nil 
}

func parseRequestLine(data []byte) (RequestLine,int ,error){

	if !bytes.Contains(data,[]byte("\r\n")) {
		return RequestLine{}, 0,nil 
	}
	bytesConsumed := bytes.Index(data, []byte("\r\n")) + 2
	
	lines := strings.Split(string(data),"\r\n")
	
	
	
	requestLine := strings.Split(lines[0]," ")
	if len(requestLine) != 3 {
		return RequestLine{}, 0,fmt.Errorf("invalid request line: %s",requestLine)
	}
	
	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(requestLine[0]) {
		return RequestLine{}, 0,fmt.Errorf("invalid method: %s",requestLine[0])
		
	}
	httpVer := strings.Split(requestLine[2],"/")
	if len(httpVer) !=2 || httpVer[0] != "HTTP" || httpVer[1] != "1.1"{
		return RequestLine{}, 0,fmt.Errorf("invalid http version: %s",requestLine[2])
	}

	
	
	method := requestLine[0]
	targetPath := requestLine[1]
	httpVersion := "1.1"
	
	return RequestLine{
		HttpVersion: httpVersion,
		RequestTarget: targetPath,
		Method: method,},bytesConsumed,nil 
}

func (r *Request) parse(data []byte) (int, error){
	totalBytesparsed := 0
	for r.HTTPReadStatus != requestStateDone{
		n,err := r.parseSingle(data[totalBytesparsed:])
		if err != nil {
			return 0, nil 
		}
		totalBytesparsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesparsed,nil 
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.HTTPReadStatus {
	case requestStateInitialized:
		rline,bytesProcessed,err := parseRequestLine(data)
		if err != nil {
			return bytesProcessed,err 
		}
		if bytesProcessed == 0 {
			return 0,nil //need more data
		}
		if bytesProcessed > 0 {
			r.RequestLine = rline 
			r.HTTPReadStatus = requestStateParsingHeaders			
			return bytesProcessed,nil
		}
		return bytesProcessed,nil 
	case requestStateParsingHeaders:
		bytesConsumed,completed,err := r.Headers.Parse(data)
		if err != nil {
			return 0, err 
		}

		if  completed{
			r.HTTPReadStatus = requestStateParsingBody			
		}
		return bytesConsumed,nil 
	case requestStateParsingBody:		
		
		bodyLength, err := r.Headers.Get("content-length")
		if err != nil {
			r.HTTPReadStatus = requestStateDone
			return len(data),nil
		}
		bodyLengthInt,err := strconv.Atoi(bodyLength)
		if err != nil {
			r.HTTPReadStatus = requestStateDone 
			return 0, err 
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		if r.bodyLengthRead > bodyLengthInt{
			return 0, fmt.Errorf("content-length too large")
		}
		if r.bodyLengthRead == bodyLengthInt{
			r.HTTPReadStatus = requestStateDone
		}
		return len(data),nil 
	case requestStateDone:
		return 0, fmt.Errorf("reading from done request")
	default:
		return 0,fmt.Errorf("error: unknown state")
		
	}
	
}