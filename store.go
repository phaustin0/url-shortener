package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateShortUrl(url Url) error
	GetUrlFromShortUrl(shortUrl string) *Url
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=urlshortener sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createUrlTable()
}

func (s *PostgresStore) createUrlTable() error {
	query := `create table if not exists url (
    shortUrl varchar(7),
    redirectUrl text 
  )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateShortUrl(url Url) error {
	query := `insert into url (shortUrl, redirectUrl) values ($1, $2)`
	_, err := s.db.Exec(query, url.ShortUrl, url.RedirectUrl)

	return err
}

func (s *PostgresStore) GetUrlFromShortUrl(shortUrl string) *Url {
	return &Url{
		ShortUrl:    shortUrl,
		RedirectUrl: "",
	}
}
