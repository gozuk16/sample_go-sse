package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sse(w http.ResponseWriter, r *http.Request) {
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				fmt.Fprintf(w, "data: %s\n\n", time.Now().Format("2006/1/2 15:04:05.000"))
				flusher.Flush()
			}
		}
	}()
	<-r.Context().Done()
	log.Println("close connection")
}

func main() {
	dir := http.Dir("./static")
	http.HandleFunc("/event", sse)
	http.Handle("/", http.FileServer(dir))
	http.ListenAndServe(":8080", nil)
}
