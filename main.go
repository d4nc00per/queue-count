package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const url string = "https://wrapapi.com/use/d4nc00per/stackoverflow/reviewQueues/0.0.2?wrapAPIKey=r7iMkqETsrMdynSjmpPFs29QRoNv1XPY"

func main() {
	q, err := getQueues()

	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		for _, r := range q {
			fmt.Printf("Queue: %v\n", r)
		}
	}
}

type queue struct {
	Type, Count string
}

type apiResponse struct {
	Success bool
	Data    struct {
		ReviewQueues []*struct {
			Queue queue
		}
	}
}

func getQueues() ([]*queue, error) {
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

	var jsonResp apiResponse

	err = json.Unmarshal(body.Bytes(), &jsonResp)

	if err != nil {
		return nil, err
	}

	if !jsonResp.Success {
		return nil, errors.New("Unsuccessful request")
	}

	queues := []*queue{}

	for _, r := range jsonResp.Data.ReviewQueues {
		queues = append(queues, &r.Queue)
	}

	return queues, nil
}
