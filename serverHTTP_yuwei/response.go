package main

import (
	"io"
	"strconv"
	"strings"
)

// Response : http response
type Response struct {
	Status    int    // "200", "404", ...
	Protocole string // "HTTP 1.1"
	Head      Header
	Body      io.Reader
	Request   *Request
}

func buildResponse(resp *Response) (respb []byte) {
	//status line
	lineS := resp.Protocole + " " + strconv.Itoa(resp.Status) + " " + StatusText(resp.Status) + "\r\n"
	respb = append(respb, lineS...)
	//fmt.Print(lineS)

	//Head
	for k, v := range resp.Head {
		lineH := k + ": " + strings.Join(v, ";") + "\r\n"
		respb = append(respb, lineH...)
		//fmt.Print(lineH)
	}
	respb = append(respb, "\r\n"...)
	/*
		//Body
		var lineB []byte
		lineB, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		respb = append(respb, lineB...)
		fmt.Print(lineB)*/
	return respb
}
