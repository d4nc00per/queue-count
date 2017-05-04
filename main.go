package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var bind string
var mongoURL string

func main() {
	mongoHost := os.Getenv("OPENSHIFT_MONGODB_DB_HOST")
	if mongoHost != "" {
		mongoURL = fmt.Sprintf(
			"%s:%s@%s:%s",
			os.Getenv("OPENSHIFT_MONGODB_DB_USERNAME"),
			os.Getenv("OPENSHIFT_MONGODB_DB_PASSWORD"),
			mongoHost,
			os.Getenv("OPENSHIFT_MONGODB_DB_PORT"))
	} else {
		mongoURL = "localhost:27017"
	}

	ip := os.Getenv("OPENSHIFT_GO_IP")
	port := os.Getenv("OPENSHIFT_GO_PORT")

	if ip != "" {
		bind = fmt.Sprintf("%s:%s", ip, port)
	} else {
		bind = "localhost:9191"
	}

	fmt.Printf("listening on %s...", bind)
	fmt.Printf("using mongo on %s...", mongoURL)

	http.HandleFunc("/update", updateQueues)

	log.Fatal(http.ListenAndServe(bind, nil))
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
