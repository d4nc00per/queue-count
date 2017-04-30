package main

import (
	"encoding/json"
	"errors"
	"time"
)

const url string = "https://wrapapi.com/use/d4nc00per/stackoverflow/reviewQueues/0.0.5?wrapAPIKey=r7iMkqETsrMdynSjmpPFs29QRoNv1XPY"

// Queue holds the name and count of a review queue
type Queue struct {
	Name, Count string
	Time        time.Time
}

type apiResponse struct {
	Success bool
	Data    struct {
		ReviewQueues []*struct {
			Q Queue
		}
	}
}

// QueueService implements methods to access the queues
type QueueService struct {
	Client *HttpClient
}

// NewQueueService creates a new instance
func NewQueueService(client *HttpClient) *QueueService {
	return &QueueService{client}
}

// GetQueues from the stack overflow api
func (that *QueueService) GetQueues() ([]*Queue, error) {
	body, err := that.Client.Get(url)

	if err != nil {
		return nil, err
	}

	var jsonResp apiResponse

	err = json.Unmarshal(body, &jsonResp)

	if err != nil {
		return nil, err
	}

	if !jsonResp.Success {
		return nil, errors.New("Unsuccessful request")
	}

	queues := []*Queue{}

	now := time.Now()

	for _, r := range jsonResp.Data.ReviewQueues {

		current := &r.Q

		current.Time = now

		queues = append(queues, current)
	}

	return queues, nil
}