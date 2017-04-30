package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/update", updateQueues)

	log.Fatal(http.ListenAndServe(":9191", nil))
}

func updateQueues(w http.ResponseWriter, r *http.Request) {
	mongo := GetDbClient()
	qs := NewQueueService(&HttpClient{})

	queues, err := qs.GetQueues()

	if err != nil {
		mongo.Log(fmt.Sprintf("Unable to retrieve queues: %v", err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	for _, q := range queues {
		err := mongo.Queues().Insert(q)
		if err != nil {
			mongo.Log(fmt.Sprintf("Error storing the queues: %v", err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	mongo.Log("Retrieved queues.")
	w.WriteHeader(http.StatusOK)
}
