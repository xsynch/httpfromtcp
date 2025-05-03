package server

import (
	"bytes"
	"fmt"
	
	"log"
	"net"
	"sync/atomic"

	"github.com/xsynch/httpfromtcp/internal/request"
	"github.com/xsynch/httpfromtcp/internal/response"
)

type Server struct{
	listener net.Listener
	closed atomic.Bool
	handler Handler
}


func Serve(port int, h Handler) (*Server, error){
	l, err := net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err != nil {
		log.Printf("error listening on port %d with error: %s",port,err)
		return nil, err 
	}
	s := Server{
		listener: l,
		handler: h,
		
		
	}
	// s.closed.Store(true)
	go s.listen()

	

	return &s,nil
}

// func HandleError(w io.Writer, he HandlerError) error {

// 	msg := he.Message
	
// 	log.Printf("Writing the status line: %d\n",he.Status)
// 	h := response.GetDefaultHeaders(len(msg))
// 	err := response.WriteStatusLine(w,response.StatusCode(he.Status))
// 	if err != nil {
		
// 		return fmt.Errorf("error writing the status line: %s",err)
// 	}
// 	log.Println("Finished writing status. Writing the headers")
// 	err = response.WriteHeaders(w,h)
// 	if err != nil {
// 		return fmt.Errorf("error writing the headers line: %s",err)
// 	}
// 	log.Println("finished writing the headers. writing the body")
// 	_, err = w.Write([]byte(msg))
// 	if err != nil {
// 		return fmt.Errorf("error writing the message line: %s",err) 
// 	}
// 	return nil 
	
// }

func (s *Server) Close() error{
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil 
}

func (s *Server) listen(){
	for {
		// if !s.closed.Load(){
		// 	log.Println("server stopped, exiting loop")
		// 	return
		// }
		conn, err := s.listener.Accept()
		if err != nil {

			if s.closed.Load(){
				log.Println("server is shutdown, not accepting new connections")
				return
			}
			log.Printf("error accepting new connection: %s",err)
			continue
		}
		
		
		go s.handle(conn)
			
		
	}
	

}

func (s *Server) handle(conn net.Conn) {
		// msg := fmt.Sprintf("HTTP/1.1 %d %s\r\n",200,"OK")
		// conn.Write([]byte(msg))
		

		defer conn.Close()
		
		
		buf := bytes.NewBuffer([]byte{})
		
		req, err := request.RequestFromReader(conn)
		if err != nil {
			
			he := &HandlerError{
				Status: 500,
				Message: err.Error(),
			}
			he.Write(conn)
			return
		}
		
		// n,err := conn.Write([]byte("Hello\r\n"))

		// log.Printf("Wrote %d bytes\n",n)
	
		handlerError := s.handler(buf,req)		
		if handlerError != nil  {
			
			handlerError.Write(conn)
			return

		} 
		b := buf.Bytes()
		defaultHeaders := response.GetDefaultHeaders(len(b))
		err = response.WriteStatusLine(conn,response.OK)
		if err != nil {
			log.Println(err)
			return
		}
	
		err = response.WriteHeaders(conn,defaultHeaders)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = conn.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
			
			
		
		// s.Close()

		// err := response.WriteStatusLine(conn,response.OK)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// h := response.GetDefaultHeaders(0)
		// err = response.WriteHeaders(conn,h)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }


		// s.Close()
}