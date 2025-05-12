package headers

import (
	"bytes"
	"fmt"
	"log"

	"regexp"
	"strings"
)


type Headers map[string]string

type State int 
const (
	initialized State = iota
	done 
)
const (
bufferSize = 8
)


type RequestHeaders struct {
	HeaderLine Headers
	HTTPReadStatus State  
}



type HeaderLine struct {
	FieldName string 
	FieldValue string 
}





func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	if len(data) == 0 {
		return 0,false,nil 
	}

	if bytes.HasPrefix(data,[]byte("\r\n")){
		return 2, true,nil 
	}

	idx := bytes.Index(data,[]byte("\r\n"))
	if idx == -1 {
		return 0,false,nil
	}

	line := data[:idx]

	colonIdx := bytes.Index(line,[]byte(":"))
	if colonIdx == -1 {
		return 2,false,fmt.Errorf("invalid header, no colon found in %s",string(line))
	}
	key := string(line[:colonIdx])
	if string(bytes.TrimRight(line[:colonIdx]," ")) != key {
		return 0, false,fmt.Errorf("invalid header, spaces before colon: %s",string(line[:colonIdx]))
	}

	key = strings.TrimSpace(key)
	val := string(bytes.TrimSpace((line[colonIdx+1:])))
	if len(val) < 1 {
		return 0,false,fmt.Errorf("header value must contain a value: %s",val)
	}
	if isValidHeader(key) {
		key = strings.ToLower(key)
		v,ok := h[key]
		if ok {			
			h[key] = fmt.Sprintf("%s, %s",v, val)
			
		} else {
			h[key] = val
		}
	} else {
		return 0,false,fmt.Errorf("key contains invalid character: %s",key)
	}
	return idx + 2,false,nil 



	

}

func isValidHeader(word string) bool {
	pattern := "^[a-zA-Z0-9\\!#\\$\\%&'\\*\\+-\\.^_`|~]*$"
	// escapedPattern := regexp.QuoteMeta(pattern)
	regex, _ := regexp.Compile(pattern)
	return regex.MatchString(word)
}

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Get(key string) (string,error) {
	val,ok := h[strings.ToLower(key)]
	if ok {
		return val,nil
	} else {
		return "",fmt.Errorf("header not found: %s",val)
	}
	

}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) OverRide(key string, val string) error {
	_,err  := h.Get(key)
	if err != nil  {
		// log.Printf("%s not found setting it to %s\n",key,err)
		h.Set(key,val)
		return nil 
	}
	log.Printf("Setting %s to %s\n",key,val)
	h[key] = val 
	return nil 

}

func (h Headers) Remove(key string) error {
	// _, err := h.Get(key)
	// if err != nil {
	// 	log.Printf("Error getting %s to remove: %s\n",key,err)
	// 	return err
	// }
	delete(h,key)
	return nil 
}