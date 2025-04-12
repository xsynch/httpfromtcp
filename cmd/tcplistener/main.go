package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/xsynch/httpfromtcp/internal/request"
)

func getLinesChannel(f io.ReadCloser) <-chan string{

	
	ch := make(chan string) 

	
	
	go func(){
		defer f.Close()
		defer close(ch)
		currentLine := ""
		
		for {
			buf := make([]byte,8,8)
			data,err := f.Read(buf)
			
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}

				if errors.Is(err,io.EOF){
					break
				}
				fmt.Printf("error reading data: %s\n",err.Error())
				return
			}


			part := strings.Split(string(buf[:data]),"\n")
			for i :=0; i < len(part)-1;i++{
				ch <- fmt.Sprintf("%s%s",currentLine,part[i])
				currentLine = ""
			}
			currentLine += part[len(part)-1]

		}
			
	}()		
	return ch

	
	
}

func main(){
	var conn net.Conn
	var err error
	listener, err := net.Listen("tcp",":42069")
	if err != nil {
		fmt.Printf("error starting listener: %s",err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn,err = listener.Accept()
		if err != nil {
			fmt.Printf("error starting listener: %s\n",err.Error())
			os.Exit(1)
		}
		fmt.Printf("connection started\n")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",req.RequestLine.Method, req.RequestLine.RequestTarget,req.RequestLine.HttpVersion)

		fmt.Printf("Connection closed\n")
		
	}


	// fileName := "messages.txt"
	// // var currentLine string

	// data,err := os.Open(fileName)
	// if err != nil {
	// 	_ = fmt.Errorf("error opening the file %s: %s",fileName,err.Error())
	// 	return 
	// }
	// defer data.Close()
	// // buf := make([]byte,8)
	
	
	// for line := range getLinesChannel(data){
	// 	fmt.Printf("read: %s\n",line)
	// }

	// for {
	// 	byteData, err := data.Read(buf)
	// 	if err != nil && err != io.EOF{
	// 		_ = fmt.Errorf("error reading data: %s",err.Error())
	// 		os.Exit(1)
	// 	}
	// 	if err == io.EOF{
	// 		// fmt.Println("End of file reached")
	// 		break
	// 	}
	// 	part := strings.Split(string(buf[:byteData]),"\n")
	// 	if len(part) > 1{
	// 		currentLine += part[0]
	// 		fmt.Printf("read: %s\n",currentLine)
	// 		currentLine = part[1]
	// 		continue
	// 	} else {
	// 		currentLine += part[0]
	// 	}
		

	// }
	
	

	
}