package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

// HTTPClient implements the http operations
type HTTPClient struct {
}

// HTTPOperations defines the http operations
type HTTPOperations interface {
	Get(url string) ([]byte, error)
}

//Get an internet resource
func (that *HTTPClient) Get(url string) ([]byte, error) {
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
