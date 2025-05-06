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

type WriterState int 

const (
	WriterStateStatusLine WriterState = iota
	WriterStateHeaders
	WriterStateBody
)

type Writer struct{
	IoWriter io.Writer

	State WriterState
	
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		IoWriter: w,
		State: WriterStateStatusLine,	
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != WriterStateStatusLine {
		return fmt.Errorf("incorrect state, should be writing the status line")
	}
	statusLine := w.GetStatusLine(statusCode)
	log.Println(string(statusLine))
	_,err := w.IoWriter.Write(statusLine)	
	if err != nil {
		return err 
	}
	w.State = WriterStateHeaders
	return nil 
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != WriterStateHeaders {
		return fmt.Errorf("incorrecct state, should be writing the headers")
	}
	for k,v := range headers {				
		log.Printf("key: %s val: %s\n",k,v)
		_, err := w.IoWriter.Write([]byte(fmt.Sprintf("%s: %s\r\n",k,v)))
		if err != nil {
			return err 
		}
	}	
	w.IoWriter.Write([]byte("\r\n"))
	w.State = WriterStateBody
	return nil 
}

func (w *Writer) WriteBody(p []byte) (int, error){
	if w.State != WriterStateBody {
		return 0,fmt.Errorf("incorrect state, should be writing the body")
	}
	n, err := w.IoWriter.Write([]byte(fmt.Sprintf("%s\n",p)))
	if err != nil {
		return 0,fmt.Errorf("error writing the body")
	}
	
	return n,nil 
}

func (w *Writer) GetStatusLine(statusCode StatusCode) []byte {
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
