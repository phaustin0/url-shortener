package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// api
// create server to listen to requests
// server does two things
//  - get request from a shortened url => redirect to url
//  - post request to create a shortened url => return shortened url

type (
	ApiFunc     func(w http.ResponseWriter, r *http.Request) error
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
)

// Server is the implementation of the server that is
// listening for requests
type Server struct {
	listenAddr string
}

type ServerError struct {
	Error string `json:"error"`
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Listen() {
	// GET request
	// POST request

	mux := handleRequests()
	log.Println("[SERVER]: running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, mux)
}

func handleRequests() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{url}", makeHandlerFunc(handleGetUrlRequest))
	mux.HandleFunc("POST /", makeHandlerFunc(handleCreateShortUrlRequest))

	return mux
}

func handleGetUrlRequest(w http.ResponseWriter, r *http.Request) error {
	url := r.PathValue("url")
	log.Println("get request for url:", url)
	return nil
}

func handleCreateShortUrlRequest(w http.ResponseWriter, r *http.Request) error {
	log.Println("post request")
	return nil
}

func makeHandlerFunc(apiFunc ApiFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := apiFunc(w, r)
		if err != nil {
			errStr := "bad request"
			WriteJSONResponse(w, http.StatusBadRequest, ServerError{Error: errStr})
			http.Error(w, errStr, http.StatusBadRequest)
		}
	}
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, v any) {
	json.NewEncoder(w).Encode(v)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
