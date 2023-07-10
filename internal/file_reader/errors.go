package file_reader

import "errors"

var (
	ErrFileDoesNotExist = errors.New("ports json file does not exist")
	ErrNoFilePath       = errors.New("filePath is required for FileReader")
)
