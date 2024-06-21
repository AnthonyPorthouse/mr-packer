package modrinth

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func ExtractOverrides(filename string, env Environment, appFs afero.Fs) error {
	info, err := appFs.Stat(filename)

	if errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if !strings.HasSuffix(filename, ".mrpack") {
		return errors.New("file is not a .mrpack")
	}

	fileReader, err := appFs.OpenFile(filename, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer fileReader.Close()

	zipReader, err := zip.NewReader(fileReader, info.Size())
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "overrides") && !strings.HasSuffix(file.Name, "/") {
			extractFile(file, appFs)
		}

		switch env {
		case Client:
			if strings.HasPrefix(file.Name, "client-overrides") && !strings.HasSuffix(file.Name, "/") {
				extractFile(file, appFs)
			}
		case Server:
			if strings.HasPrefix(file.Name, "client-overrides") && !strings.HasSuffix(file.Name, "/") {
				extractFile(file, appFs)
			}
		}
	}

	return nil
}

func extractFile(file *zip.File, appFs afero.Fs) error {

	targetFile := string(file.Name[strings.IndexRune(file.Name, '/')+1:])

	log.Printf("Extracting %s to %s", file.Name, targetFile)

	_, err := appFs.Stat(targetFile)

	if errors.Is(err, fs.ErrNotExist) {
		dir := filepath.Dir(targetFile)

		if _, err := appFs.Stat(dir); errors.Is(err, fs.ErrNotExist) {
			appFs.MkdirAll(dir, 0755)
		}
	}

	originalFile, err := file.Open()
	if err != nil {
		return err
	}
	defer originalFile.Close()

	newFile, err := appFs.Create(targetFile)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err = io.Copy(newFile, originalFile); err != nil {
		return err
	}

	return nil
}
