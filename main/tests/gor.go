package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync/atomic"
)

var recording int32

func isRecording() bool {
	return atomic.LoadInt32(&recording) != 0
}

func setRecording(shouldRecord bool) {
	if shouldRecord {
		atomic.StoreInt32(&recording, 1)
	} else {
		atomic.StoreInt32(&recording, 0)
	}
}

type SwitchHandler struct {
	mux http.Handler
}

func (s *SwitchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isRecording() {
		fmt.Printf("Switch Handler is Recording\n")
		s.mux.ServeHTTP(w, r)
		return
	}

	fmt.Printf("Switch Handler is NOT Recording\n")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "NOT Recording\n")

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/success/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Recording\n")
	})

	handler := &SwitchHandler{mux: router}

	setRecording(true)

	http.Handle("/", handler)

	http.ListenAndServe(":8080", nil)
}
