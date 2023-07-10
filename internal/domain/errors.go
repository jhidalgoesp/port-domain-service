package domain

import "errors"

var (
	ErrNoFileRepository = errors.New("fileRepository is required")
	ErrNoPortRepository = errors.New("portRepositoy is required")
)
