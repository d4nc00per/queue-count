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
	setupMongoURL()
	setupAppURL()

	fmt.Printf("listening on %s...\n", bind)
	fmt.Printf("using mongo on %s...\n", mongoURL)

	http.HandleFunc("/update", updateQueues)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/index.html")
		log.Print("Index returned.")
	})

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

	log.Print("Queues updated.")
	w.WriteHeader(http.StatusOK)
}

func setupAppURL() {
	ip := os.Getenv("OPENSHIFT_GO_IP")
	port := os.Getenv("OPENSHIFT_GO_PORT")

	if ip != "" {
		bind = fmt.Sprintf("%s:%s", ip, port)
	} else {
		bind = "localhost:9191"
	}
}

func setupMongoURL() {
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
}
