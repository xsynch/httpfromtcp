package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main(){
	fileName := "messages.txt"
	var currentLine string
	data,err := os.Open(fileName)
	if err != nil {
		_ = fmt.Errorf("error opening the file %s: %s",fileName,err.Error())
		return 
	}
	defer data.Close()
	buf := make([]byte,8)

	for {
		byteData, err := data.Read(buf)
		if err != nil && err != io.EOF{
			_ = fmt.Errorf("error reading data: %s",err.Error())
			os.Exit(1)
		}
		if err == io.EOF{
			// fmt.Println("End of file reached")
			break
		}
		part := strings.Split(string(buf[:byteData]),"\n")
		if len(part) > 1{
			currentLine += part[0]
			fmt.Printf("read: %s\n",currentLine)
			currentLine = part[1]
			continue
		} else {
			currentLine += part[0]
		}
		

	}
	// fmt.Printf("read: %s\n",currentLine)

	
}