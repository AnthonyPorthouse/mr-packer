package modrinth

import (
	"bytes"
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

func DownloadFiles(manifest *Manifest, environment Environment, appFs afero.Fs) error {

	errs, ctx := errgroup.WithContext(context.Background())

	for _, file := range manifest.Files {
		thisFile := *file

		errs.Go(func() error {
			return downloadFile(ctx, thisFile, environment, appFs)
		})
	}

	resultErrs := errs.Wait()

	return resultErrs
}

func validateDownloadPath(path string) error {

	normalizedPath := filepath.Clean(path)

	if strings.HasPrefix(normalizedPath, "../") || strings.HasPrefix(normalizedPath, "/") {
		return errors.New("cannot navigate below instance directory")
	}

	return nil
}

func downloadFile(ctx context.Context, file ManifestFile, environment Environment, appFs afero.Fs) error {
	switch environment {
	case Client:
		if file.Env.Client == Unsupported {
			return nil
		}
	case Server:
		if file.Env.Server == Unsupported {
			return nil
		}
	}

	log.Printf("Downloading %s", file.Path)

	if err := validateDownloadPath(file.Path); err != nil {
		return err
	}

	stat, err := appFs.Stat(file.Path)

	if errors.Is(err, fs.ErrNotExist) {
		dir := filepath.Dir(file.Path)

		if _, err := appFs.Stat(dir); errors.Is(err, fs.ErrNotExist) {
			log.Printf("Making Directory %s", dir)
			appFs.MkdirAll(dir, 0755)
		}
	} else if err != nil {
		return err
	}

	if stat != nil && stat.Size() == file.FileSize {
		log.Println("File already exists and valid, skipping")

		return nil
	}

	newFile, err := appFs.Create(file.Path)
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
				return false, fmt.Errorf("%s hash does not match %s", hash, file.Hashes.Sha512)
			}

			return true, nil
		}()

		errs = append(errs, err)

		if success {
			break
		}
	}

	return errors.Join(errs...)

}
