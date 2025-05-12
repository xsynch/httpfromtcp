package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	var streamInfo string 

	log.Printf("Attempting to target: %s\n",req.RequestLine.RequestTarget)
	h := headers.Headers{}
	binCheck := strings.HasPrefix(req.RequestLine.RequestTarget,"/httpbin")
	if binCheck {
		streamInfo = strings.TrimPrefix(req.RequestLine.RequestTarget,"/httpbin/")
		req.RequestLine.RequestTarget = "/httpbin"
	}
	
	

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.BadRequest
		msg ,err = (getBody(response.BadRequest))
		if err != nil {
			log.Println("Error with getting the request body")
			return 
		}
	case "/video":
		status = response.OK
		video, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			log.Printf("Error reading the video file: %s\n",err)
			return 
		}
		
		err = w.WriteStatusLine(status)
		if err != nil {
			log.Printf("Error writing the status line: %s\n",err)
			return 
		}

		h = response.GetDefaultHeaders(len(video))
		h.OverRide("content-type","video/mp4")
		err = w.WriteHeaders(h)		
		if err != nil {
			log.Printf("error writing the headers line: %s\n",err)
			return 
		}

		_,err = w.IoWriter.Write(video)
		if err != nil {
			log.Printf("Error playing the video: %s\n",err)
			return 
		}
		return 


	case "/myproblem":
		status = response.InternalServerError
		msg ,err = getBody(status)
		if err != nil{
			log.Println("Error with getting the request body")
			return
		}
	case "/httpbin":
		
		buf := make([]byte,1024)
		totalWrittenBytes := 0
		bodyData := ""
		
		url := fmt.Sprintf("https://httpbin.org/%s",streamInfo)
		resp,err := http.Get(url)
		if err != nil {
			log.Printf("Error getting information: %s",err)
			return 
		}
		h = response.GetDefaultHeaders(len(buf))
		err = h.Remove("content-length")
		if err != nil {
			log.Printf("Error removing content-length: %s\n",err)
			return
		}
		h.Set("transfer-encoding","chunked")
		err = w.WriteStatusLine(response.OK)
		if err != nil {
			log.Printf("Error writing the status line: %s\n",err)
			return 
		}
		
		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("Error writing headers: %s\n",err)
			return 
		}
		defer resp.Body.Close()
		for {
			
			n,err := resp.Body.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF){
					_,err = w.WriteChunkedBodyDone()
					if err != nil {
						log.Printf("Error writing end of chunk: %s\n",err)
						return 
					}
					break 					
				}
				log.Printf("Error reading chunk: %d %s\n",n,err)
				return  
			}
			bodyData += string(buf[:n])
			b,err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("Error writing the body chunk: %s\n",err)
				return
			}
			totalWrittenBytes += b 
		}
		bodyHash := sha256.Sum256([]byte(bodyData))
		trailerHeaders := headers.NewHeaders()
		trailerHeaders.Remove("content-length")
		trailerHeaders.Set("X-Content-SHA256",fmt.Sprintf("%x",bodyHash))
		trailerHeaders.Set("X-Content-Length",fmt.Sprintf("%d",len(bodyData)))
		w.WriteTrailers(trailerHeaders)
		return 

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