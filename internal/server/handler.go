package server

import (
	"io"
	"log"

	"github.com/xsynch/httpfromtcp/internal/request"
	"github.com/xsynch/httpfromtcp/internal/response"
)

type HandlerError struct {
	Status int 
	Message string 
}

type Handler func(w io.Writer, req *request.Request) *HandlerError


func (he HandlerError) Write(w io.Writer) {
	// log.Println("writing status line")	
	response.WriteStatusLine(w, response.StatusCode(he.Status))
	
	messageBytes := []byte(he.Message)
	// log.Println("completed status line, writing headers")
	headers := response.GetDefaultHeaders(len(messageBytes))	
	err := response.WriteHeaders(w, headers)
	if err != nil {
		log.Printf("Error writing headers: %s\n",err)
		return
	}
	_,err = w.Write(messageBytes)
	if err != nil{
		log.Printf("Error writing body: %s\n",err)
		return 
	}
}