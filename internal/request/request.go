package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/xsynch/httpfromtcp/internal/headers"
	
)

type State int 
const (
	initialized State = iota
	requestStateParsingHeaders 
	parsingBody
	done 
)
const (
bufferSize = 8
)


type Request struct {
	RequestLine RequestLine
	HTTPReadStatus State 
	Headers headers.Headers
	Body []byte
	
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
		HTTPReadStatus: initialized,
		Headers: map[string]string{},
	}
	
	
	for req.HTTPReadStatus == initialized{
		
		if readToindex == len(buf){
			newBuf := make([]byte,len(buf) * 2)
			_ = copy(newBuf,buf)
			buf = newBuf

		}
		b,err := reader.Read(buf[readToindex:])
		if err != nil && err != io.EOF{
			return nil, err				
		}
		readToindex += b
		
		bytesConsumed,err := req.parse(buf[:readToindex])
		if err != nil {
			return nil, err 
		}
		if bytesConsumed > 0 {
			copy(buf,buf[bytesConsumed:readToindex])
			readToindex -= bytesConsumed
		}

		if err == io.EOF{
			req.HTTPReadStatus = done
			break
		}
				
	}
	 for req.HTTPReadStatus == requestStateParsingHeaders{

		if readToindex == len(buf){
			newBuf := make([]byte,len(buf) * 2)
			_ = copy(newBuf,buf)
			buf = newBuf

		}
		b,err := reader.Read(buf[readToindex:])
		if err != nil && err != io.EOF{
			return nil, err				
		}
		readToindex += b
		
		bytesConsumed,completed,err := req.Headers.Parse(buf[:readToindex])
		if err != nil {
			return nil, err 
		}
		if bytesConsumed > 0 {
			copy(buf,buf[bytesConsumed:readToindex])
			readToindex -= bytesConsumed
		}

		if  completed{
			req.HTTPReadStatus = parsingBody
			// break
		}
	 }
	 for req.HTTPReadStatus == parsingBody {
		totalBodyBytes := 0
		bodyLength, err := req.Headers.Get("content-length")
		if err != nil {
			req.HTTPReadStatus = done
			break
		}
		bodyLengthInt,err := strconv.Atoi(bodyLength)
		if err != nil {
			req.HTTPReadStatus = done 
			return nil, err 
		}
		//readtoindex already contains the first two bytes of the body
		totalBodyBytes = readToindex
		req.Body = append(req.Body, buf[:readToindex]...)
		for {
			if readToindex == len(buf){
				newBuf := make([]byte,len(buf) * 2)
				_ = copy(newBuf,buf)
				buf = newBuf
	
			}
			// fmt.Printf("Read to index val: %d\n",readToindex)
			n, err := reader.Read(buf[readToindex:])			
			if err == io.EOF{
				
				if totalBodyBytes > bodyLengthInt{
					req.HTTPReadStatus = done 
					return nil, fmt.Errorf("body contents greater than content length provided")
				}
				
				if totalBodyBytes < bodyLengthInt {
					req.HTTPReadStatus = done 
					return nil, fmt.Errorf("body contents less than content length provided")
				}
				if totalBodyBytes == bodyLengthInt {
					req.HTTPReadStatus = done
					fmt.Printf("consumed %d of the body data\n",totalBodyBytes)
					break
				}
				
				
			} else if err != nil {				
				return nil,err 
			}
			
			if n > 0 {
												
				req.Body = append(req.Body, buf[readToindex:readToindex+n]...)
				readToindex += n
				totalBodyBytes += n 

			} else {
				req.HTTPReadStatus = done
				break
				// return nil,fmt.Errorf("error reading buffer infomration")
			}
		}		
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

func (r *Request) parse(data []byte) (int, error) {
	if r.HTTPReadStatus == initialized{
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
	}

	if r.HTTPReadStatus == done {
		return 0,fmt.Errorf("error: trying to read data in a done state")
	}

	return 0,fmt.Errorf("error: unknown state")
}