package main

import (
	"net/http"
	"sse/broker"
)

type Widget struct {
	Id    int
	Attrs []string
}

func main() {
	r := http.NewServeMux()

	broker := broker.New()

	r.HandleFunc("/sse", broker.Stream)

	http.ListenAndServe(":3333", r)
}
