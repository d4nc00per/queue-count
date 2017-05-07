package main

import (
	"encoding/json"
	"errors"
	"time"
)

const url string = "https://wrapapi.com/use/d4nc00per/stackoverflow/reviewQueues/0.0.5?wrapAPIKey=r7iMkqETsrMdynSjmpPFs29QRoNv1XPY"

// Queue holds the name and count of a review queue
type Queue struct {
	Name  string    `bson:"name"`
	Count string    `bson:"count"`
	Time  time.Time `bson:"time"`
}

type apiResponse struct {
	Success bool
	Data    struct {
		ReviewQueues []*struct {
			Q Queue `json:"queue"`
		} `json:"review_queues"`
	}
}

// QueueService implements methods to access the queues
type QueueService struct {
	Client HTTPOperations
}

// NewQueueService creates a new instance
func NewQueueService(client HTTPOperations) *QueueService {
	return &QueueService{client}
}

// GetQueues from the stack overflow api
func (that *QueueService) GetQueues() ([]*Queue, error) {
	body, err := that.Client.Get(url)

	if err != nil {
		return nil, err
	}

	var jsonResp apiResponse

	// log.Print(string(body))

	err = json.Unmarshal(body, &jsonResp)

	if err != nil {
		return nil, err
	}

	// log.Printf("%v", jsonResp)

	if !jsonResp.Success {
		return nil, errors.New("Unsuccessful request")
	}
	queues := []*Queue{}

	now := time.Now().UTC()

	for _, r := range jsonResp.Data.ReviewQueues {

		current := &r.Q

		current.Time = now

		queues = append(queues, current)
	}

	return queues, nil
}
