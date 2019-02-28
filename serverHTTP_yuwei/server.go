package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func closeConnection(conn net.Conn) {
	conn.Close()
}

func newConnection() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer closeConnection(conn)
	//read from connection
	reqb := make([]byte, 1024)
	_, err := conn.Read(reqb)
	if err != nil {
		log.Fatal(err)
	}
	//parse request
	req, err := parseRequest(reqb)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s %s\n", req.getMethod(), req.getURL(), req.getProtocole())
	fmt.Println(req.Head)

	//write to connection
	resp := new(Response)
	switch req.getMethod() {
	case "GET":
		conn.Write(doGET(req, resp))
	case "POST":
		conn.Write(doPOST(req, resp))
	}

	//conn.Write(buildResponse(resp))
	//fmt.Println(buildResponseStr(resp))
}

func doGET(req *Request, resp *Response) (respb []byte) {
	resp.Status = 200 //404
	resp.Protocole = "HTTP/1.1"
	resp.Request = req

	file, err := os.Open("post.html")
	if err != nil {
		log.Fatal(err)
	}
	resp.Body = file
	defer file.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Head = Header{
		"Content-Type":   []string{"text/html"},
		"Content-Length": []string{strconv.Itoa(len(body))},
	}
	respb = append(respb, buildResponse(resp)...)
	respb = append(respb, body...)
	return respb
}

func doPOST(req *Request, resp *Response) (respb []byte) {
	resp.Status = 200 //404
	resp.Protocole = "HTTP/1.1"
	resp.Request = req

	reqBody := parsePost(req) //map[string][]string
	var body string
	for k, v := range reqBody {
		body = body + k + ": " + strings.Join(v, ",") + "\r\n"
	}
	resp.Body = strings.NewReader(body)

	resp.Head = Header{
		"Content-Type":   []string{"text/html"},
		"Content-Length": []string{strconv.Itoa(len(body))},
	}
	respb = append(respb, buildResponse(resp)...)
	respb = append(respb, body...)
	return respb
}
