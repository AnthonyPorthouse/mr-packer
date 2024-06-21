package internal

import (
	"archive/zip"
	"errors"
	"log"
	"os"
	"strings"
)

func ValidateFile(filename string) (bool, error) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, err
	}

	if !strings.HasSuffix(filename, ".mrpack") {
		return false, errors.New("file is not a .mrpack")
	}

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		return false, err
	}

	defer zipReader.Close()
	for _, file := range zipReader.File {
		log.Println(file.Name)
		if file.Name == "modrinth.index.json" {
			return true, nil
		}
	}

	return false, nil
}
