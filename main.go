package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2/bson"
)

var bind string
var mongoURL string

func main() {
	setupMongoURL()
	setupAppURL()

	fmt.Printf("listening on %s...\n", bind)
	fmt.Printf("using mongo on %s...\n", mongoURL)

	http.HandleFunc("/update", updateQueues)
	http.HandleFunc("/data", getQueues)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/index.html")
		log.Print("Index returned.")
	})

	http.HandleFunc("/d3", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/d3.html")
		log.Print("D3 returned.")
	})

	scripts := http.FileServer(http.Dir("scripts"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))

	d3 := http.FileServer(http.Dir("d3"))
	http.Handle("/d3/", http.StripPrefix("/d3/", d3))

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

func getQueues(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving the data")

	var queues []*Queue
	mongo := GetDbClient(mongoURL)
	err := mongo.Queues().Find(bson.M{}).All(&queues)

	log.Printf("Queue count:%d", len(queues))
	if err != nil {
		mongo.Log(fmt.Sprintf("Error querying for the queues: %v", err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	data, err := json.Marshal(queues)

	if err != nil {
		mongo.Log(fmt.Sprintf("Error marshalling the queues: %v", err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	_, err = w.Write(data)

	if err != nil {
		mongo.Log(fmt.Sprintf("Error marshalling the queues: %v", err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	log.Println("Returned the data")
	mongo.Log("Returned the data")
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
