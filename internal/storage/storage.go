/*
определяем общие ошибки для storage
и запросы к БД
*/
package storage

import "errors"

const (
	CreateTableURL = `
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`
	InsertIntoURL = `INSERT INTO url(url, alias) VALUES(?, ?);`
	DeleteFromURL = `DELETE FROM url WHERE alias == ?;`
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)
