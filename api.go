package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// api
// create server to listen to requests
// server does two things
//  - get request from a shortened url => redirect to url
//  - post request to create a shortened url => return shortened url

const shortUrlLength int = 7

type (
	ApiFunc     func(w http.ResponseWriter, r *http.Request) (int, ServerError)
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
)

type CreateShortUrlRequestBody struct {
	Url string `json:"url"`
}

type Url struct {
	ShortUrl    string `json:"shortUrl"`
	RedirectUrl string `json:"redirectUrl"`
}

func NewUrl(shortUrl, redirectUrl string) *Url {
	return &Url{
		ShortUrl:    shortUrl,
		RedirectUrl: redirectUrl,
	}
}

// Server is the implementation of the server that is
// listening for requests
type Server struct {
	listenAddr string
	store      Storage
}

type ServerError struct {
	ErrStr string `json:"error"`
}

var NoServerError = ServerError{
	ErrStr: "",
}

func NewServer(listenAddr string, store Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Listen() {
	// GET request
	// POST request

	mux := s.handleRequests()
	log.Println("[SERVER]: running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, mux)
}

func (s *Server) handleRequests() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{url}", makeHandlerFunc(s.handleGetUrlRequest))
	mux.HandleFunc("POST /", makeHandlerFunc(s.handleCreateShortUrlRequest))

	return mux
}

func (s *Server) handleGetUrlRequest(w http.ResponseWriter, r *http.Request) (int, ServerError) {
	shortUrl := r.PathValue("url")

	url, err := s.store.GetUrlFromShortUrl(shortUrl)
	if err != nil {
		errStr := fmt.Sprint("unable to read from database: %s", err.Error())
		return http.StatusInternalServerError, ServerError{ErrStr: errStr}
	}

	http.Redirect(w, r, url.RedirectUrl, http.StatusSeeOther)

	return http.StatusOK, NoServerError
}

func (s *Server) handleCreateShortUrlRequest(
	w http.ResponseWriter,
	r *http.Request,
) (int, ServerError) {
	body := new(CreateShortUrlRequestBody)
	json.NewDecoder(r.Body).Decode(body)

	if body.Url == "" {
		return http.StatusBadRequest, ServerError{ErrStr: "invalid request"}
	}

	shortUrl := getShortUrl(shortUrlLength)
	url := NewUrl(shortUrl, body.Url)

	err := s.store.CreateShortUrl(*url)
	if err != nil {
		errStr := fmt.Sprint("unable to write to database: %s", err.Error())
		return http.StatusInternalServerError, ServerError{ErrStr: errStr}
	}

	WriteJSONResponse(w, http.StatusOK, url)

	return http.StatusOK, NoServerError
}

func makeHandlerFunc(apiFunc ApiFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := apiFunc(w, r)
		if err != NoServerError {
			WriteJSONResponse(w, statusCode, err)
		}
	}
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func getShortUrl(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
