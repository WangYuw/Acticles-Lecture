package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/textproto"
	"net/url"
	"strings"
)

// Header : head of http req or resp
type Header map[string][]string

//Request : http request
type Request struct {
	Method    string // "GET", "POST", "PUT", "DELETE", ...
	URL       string
	Protocole string // "HTTP 1.1"
	Head      Header
	Body      io.Reader
}

func (req *Request) getMethod() string {
	return req.Method
}

func (req *Request) getProtocole() string {
	return req.Protocole
}

func (req *Request) getURL() string {
	return req.URL
}

//"GET...", head, body
func parseRequest(bytearr []byte) (req *Request, err error) {
	req = new(Request)
	bReader := bufio.NewReader(bytes.NewBuffer(bytearr))
	tpReader := textproto.NewReader(bReader)
	//first line
	var s string
	if s, err = tpReader.ReadLine(); err != nil {
		return nil, err
	}
	req.Method, req.URL, req.Protocole = parseRequestLine(s)
	//head
	head, err := tpReader.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	req.Head = Header(head)
	//body
	body, err := tpReader.ReadLineBytes()
	req.Body = bytes.NewBuffer(body)
	//tmp, _ := ioutil.ReadAll(req.Body)
	//fmt.Println(string(tmp))
	return req, err
}

//parse "GET /users HTTP/1.1"
func parseRequestLine(line string) (method, reqURL, protocole string) {
	strs := strings.Split(line, " ")
	if len(strs) != 3 {
		return
	}
	method = strs[0]
	reqURL = strs[1]
	protocole = strs[2]
	return method, reqURL, protocole
}

//para1=value1&para2=value2 -> type Values map[string][]string
func parsePost(req *Request) url.Values {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	body, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	return body
}
