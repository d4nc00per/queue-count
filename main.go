package main

func main() {

	qs := NewQueueService(&HttpClient{})

	qs.GetQueues()
}
