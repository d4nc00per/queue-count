package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/update", updateQueues)

	log.Fatal(http.ListenAndServe(":9191", nil))
}

func updateQueues(w http.ResponseWriter, r *http.Request) {

	qs := NewQueueService(&HttpClient{})

	qs.GetQueues()

	log.Printf("Retrieved queues.")
}
