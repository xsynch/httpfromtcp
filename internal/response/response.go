package response

import (
	"fmt"
	"io"

	"github.com/xsynch/httpfromtcp/internal/headers"
)

type StatusCode int 

const (
	OK StatusCode = 200
	BadRequest = 400
	InternalServerError = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case BadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case InternalServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		w.Write([]byte("HTTP/1.1 200\r\n"))
	}
	return nil 
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	
	h := headers.Headers{
		"Content-Length": fmt.Sprintf("%d",contentLen),
		"Connection": "close",
		"Content-Type": "text/plain",

	}
	return h 
}


func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key,val := range headers{
		w.Write([]byte(fmt.Sprintf("%s: %s\r\n",key,val)))
	}
	w.Write([]byte("\r\n"))
	return nil 
}