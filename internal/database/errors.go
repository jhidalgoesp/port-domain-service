package database

import "errors"

var (
	ErrPortNotFound = errors.New("port record not found in database")
)
