package modrinth

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
)

func ValidateFile(filename string, appFs afero.Fs) (*Manifest, error) {

	info, err := appFs.Stat(filename)

	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	if !strings.HasSuffix(filename, ".mrpack") {
		return nil, errors.New("file is not a .mrpack")
	}

	fileReader, err := appFs.OpenFile(filename, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	zipReader, err := zip.NewReader(fileReader, info.Size())
	if err != nil {
		return nil, err
	}

	for _, file := range zipReader.File {
		log.Println(file.Name)
		if file.Name != "modrinth.index.json" {
			continue
		}

		manifestFile, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer manifestFile.Close()

		return parseManifest(&manifestFile)
	}

	return nil, nil
}

func parseManifest(file *io.ReadCloser) (*Manifest, error) {
	data, err := io.ReadAll(*file)
	if err != nil {
		return nil, err
	}

	var manifest Manifest

	if err = json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, err
}
