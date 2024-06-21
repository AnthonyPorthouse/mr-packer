package internal

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"log"
	"net/http"
)

func DownloadFiles(manifest *Manifest, environment Environment, fs afero.Fs) error {
	for _, file := range manifest.Files {
		err := downloadFile(file, environment, fs)
		if err != nil {
			return err
		}
	}

	log.Printf("%+v", fs)

	return nil
}

func downloadFile(file *ManifestFile, environment Environment, fs afero.Fs) error {
	log.Printf("Downloading %s", file.Path)

	newFile, err := fs.Create(file.Path)
	if err != nil {
		return err
	}
	defer newFile.Close()

	errs := make([]error, len(file.Downloads))

	for _, url := range file.Downloads {
		success, err := func() (bool, error) {
			log.Printf("Downloading from %s", url)
			resp, err := http.Get(url)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return false, errors.New(resp.Status)
			}

			var buf bytes.Buffer
			tee := io.TeeReader(resp.Body, &buf)

			_, err = io.Copy(newFile, tee)
			if err != nil {
				return false, err
			}

			h := sha512.New()
			if _, err = io.Copy(h, &buf); err != nil {
				return false, err
			}

			hash := fmt.Sprintf("%x", h.Sum(nil))

			if hash != file.Hashes.Sha512 {
				return false, errors.New(fmt.Sprintf("%s hash does not match %s", hash, file.Hashes.Sha512))
			}

			log.Printf("Hash %s matches expected %s\n", hash, file.Hashes.Sha512)

			return true, nil
		}()

		errs = append(errs, err)

		if success == true {
			break
		}
	}

	return errors.Join(errs...)

}
