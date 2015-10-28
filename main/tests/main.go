package main

import "net/http"

type HasHandleFunc interface {
	HandleFunc(pattern string, handler func(w http.ResponseWriter, req *http.Request))
}
type Handler struct {
	http.HandlerFunc
	Enabled bool
}
type Handlers map[string]*Handler

func (h Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if handler, ok := h[path]; ok && handler.Enabled {
		//handler.ServeHTTP(w, r)
		w.Write([]byte("ACTIVATED"))
		handler.Enabled = false
		handler.ServeHTTP(w, r)
	} else {
		//http.Error(w, "Not Found", http.StatusNotFound)
		if path == "/" {
			w.Write([]byte("ACTIVATED"))
		} else if path == "/plugin1" {
			w.Write([]byte("PLUG ACTIVATED"))
		}
	}
}

func (h Handlers) HandleFunc(mux HasHandleFunc, pattern string, handler http.HandlerFunc) {
	h[pattern] = &Handler{handler, true}
	mux.HandleFunc(pattern, h.ServeHTTP)
}

func main() {
	mux := http.NewServeMux()
	handlers := Handlers{}
	handlers.HandleFunc(mux, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this will show once"))
		handlers["/"].Enabled = false
	})

	http.Handle("/", mux)
	http.ListenAndServe(":9020", nil)
}
