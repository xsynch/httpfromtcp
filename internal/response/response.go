package response

import (
	"fmt"
	"io"
	"log"

	"github.com/xsynch/httpfromtcp/internal/headers"
)

type StatusCode int 

const (
	OK StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

// func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
// 	var reasonPhrase string 
// 	switch statusCode {
// 	case OK:
// 		reasonPhrase = "OK"
// 	case BadRequest:
// 		reasonPhrase = "Bad Request"
// 	case InternalServerError:
// 		reasonPhrase = "Internal Server Error"
// 	default:
// 		reasonPhrase = ""
// 	}	
// 	_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n",statusCode,reasonPhrase)))


// 	return err 
// }

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length",fmt.Sprintf("%d",contentLen))
	h.Set("Connection","close")
	h.Set("Content-Type","text/plain")
	
	// h := headers.Headers{
	// 	"Content-Length": fmt.Sprintf("%d",contentLen),
	// 	"Connection": "close",
	// 	"Content-Type": "text/plain",

	// }
	return h 
}


func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case OK:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case InternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	return fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}


// func WriteHeaders(w io.Writer, headers headers.Headers) error {
// 	for key,val := range headers{
// 		log.Printf("Writing header: %s : %s",key,val)
		
// 		_,err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n",key,val)))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	_, err := w.Write([]byte("\r\n"))

// 	return err 
// }

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	
	for k, v := range headers {
		log.Printf("Writing header: %s: %s",k,v)
		_, err := w.Write(fmt.Appendf(nil,"%s: %s\r\n", k, v))
		if err != nil {
			return fmt.Errorf("error writing another header: %s",err)
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
