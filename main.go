package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal("[ERROR]: unable to connect to database")
	}

	if err := store.Init(); err != nil {
		log.Fatalf("[ERROR]: unable to create database tables: %s", err.Error())
	}

	listenAddr := ":8000"
	s := NewServer(listenAddr, store)
	s.Listen()
}
