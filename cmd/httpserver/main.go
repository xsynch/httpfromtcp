package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"io"

	"github.com/xsynch/httpfromtcp/internal/server"
	
	"github.com/xsynch/httpfromtcp/internal/request"
)

const port = 42069

func main() {
	
	
	server, err := server.Serve(port,handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerFunc (w io.Writer, req *request.Request) *server.HandlerError{
	
	log.Printf("Attempting to target: %s\n",req.RequestLine.RequestTarget)

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			Status: 400,
			Message: "Your problem is not my problem\n",
			
		}
	case "/myproblem":
		return &server.HandlerError{
			Status: 500,
			Message: "Woopsie, my bad\n",
		}
	default:
		_,err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			return &server.HandlerError{
				Status: 500,
				Message: "Error Writing the response\n",
			}
		}
	}

	return nil 
	
}