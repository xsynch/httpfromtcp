package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


func main(){
	hostNameandPort := "localhost:42069"

	addr,err := net.ResolveUDPAddr("udp",hostNameandPort)
	if err != nil {
		fmt.Printf("error resolving url: %s\n",err.Error())
		os.Exit(1)
	}

	udpConn,err := net.DialUDP("udp",nil,addr)
	if err != nil {
		fmt.Printf("error dialing udp connection: %s\n",err.Error())
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input,err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %s\n",err.Error())
		}
		_,err = udpConn.Write([]byte(input))
		if err != nil {
			fmt.Printf("error writing to the connection: %s\n",err.Error())
		}

	}
}


