package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

type HttpClient struct {
}

//Get an internet resource
func (that *HttpClient) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	writer := bufio.NewWriter(&body)

	defer resp.Body.Close()
	_, err = io.Copy(writer, resp.Body)

	if err != nil {
		return nil, err
	}

	return body.Bytes(), err
}
