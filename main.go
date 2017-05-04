package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var mongoURL string

func main() {

	flag.StringVar(&mongoURL, "mongo", "localhost:27017", "The url for the mongo database")

	http.HandleFunc("/update", updateQueues)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func updateQueues(w http.ResponseWriter, r *http.Request) {
	mongo := GetDbClient(mongoURL)
	qs := NewQueueService(&HTTPClient{})

	queues, err := qs.GetQueues()

	if err != nil {
		mongo.Log(fmt.Sprintf("Unable to retrieve queues: %v", err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if len(queues) == 0 {
		mongo.Log("No queues")
	}

	for _, q := range queues {
		err := mongo.Queues().Insert(q)
		if err != nil {
			mongo.Log(fmt.Sprintf("Error storing the queues: %v", err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	log.Print("Request completed.")
	mongo.Log("Request completed.")
	w.WriteHeader(http.StatusOK)
}
