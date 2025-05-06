package main

import (

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xsynch/httpfromtcp/internal/headers"
	"github.com/xsynch/httpfromtcp/internal/response"
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

func handlerFunc (w *response.Writer, req *request.Request) {
	
	var msg string
	var err error
	var status response.StatusCode

	log.Printf("Attempting to target: %s\n",req.RequestLine.RequestTarget)
	h := headers.Headers{}
	

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.BadRequest
		msg ,err = (getBody(response.BadRequest))
		if err != nil {
			log.Println("Error with getting the request body")
			return 
		}
		



	case "/myproblem":
		status = response.InternalServerError
		msg ,err = getBody(status)
		if err != nil{
			log.Println("Error with getting the request body")
			return
		}
		// h := response.GetDefaultHeaders(len(msg))
		
		// w.WriteHeaders(h)		

		// w.WriteBody([]byte(msg))

	default:
		status = response.OK
		msg,err = getBody(status)
		if err != nil {
			log.Println("Error with getting the request body")
			return
		}
		 
	}

	err = w.WriteStatusLine(status)
	if err != nil {
		log.Printf("error writing the status line: %s\n",err)
	}
	
	h = response.GetDefaultHeaders(len(msg))
	h.OverRide("content-type","text/html")
	
	err = w.WriteHeaders(h)		
	if err != nil {
		log.Printf("error writing the headers line: %s\n",err)
	}
	

	_, err = w.WriteBody([]byte(msg))	
	if err != nil {
		log.Printf("error writing the body: %s\n",err)
	}

	
	
}

func getBody(r response.StatusCode) (string,error){
	switch r {
	case response.BadRequest:
		resp := `<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
		</html>`
		return resp,nil
	case response.InternalServerError:
		resp := `<html>
				<head>
					<title>500 Internal Server Error</title>
				</head>
				<body>
					<h1>Internal Server Error</h1>
					<p>Okay, you know what? This one is on me.</p>
				</body>
				</html>` 
				return resp,nil
	default:
		resp :=  `<html>
				<head>
					<title>200 OK</title>
				</head>
				<body>
					<h1>Success!</h1>
					<p>Your request was an absolute banger.</p>
				</body>
				</html>`
		return resp,nil 
		
	}
}