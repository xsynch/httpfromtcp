package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct{
	listener net.Listener
	state atomic.Bool
}


func Serve(port int) (*Server, error){
	l, err := net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err != nil {
		log.Fatalf("error listening on port %d with error: %s",port,err)
	}
	s := Server{
		listener: l,
		
		
	}
	s.state.Store(true)
	go s.listen()
	return &s,nil
}

func (s *Server) Close() error{
	s.state.Store(false)
	err := s.listener.Close()
	if err != nil {
		return err 
	}
	return nil 
}

func (s *Server) listen(){
	for {
		if !s.state.Load(){
			log.Println("server stopped, exiting loop")
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {

			if !s.state.Load(){
				log.Println("server is shutdown, not accepting new connections")
				return
			}
			log.Printf("error accepting new connection: %s",err)
			continue
		}
		
		go func(c net.Conn) {
			s.handle(c)
			
		}(conn)
	}
	

}

func (s *Server) handle(conn net.Conn) {
		defer conn.Close()
		_,err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\r\n"))
		if err != nil {
			log.Println("error writing data to the connection")			
		}
		
		s.Close()
}