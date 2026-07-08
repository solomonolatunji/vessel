package utils

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// CreateTarContext archives the directory at sourceDir into an in-memory tarball suitable for Docker daemon builds.
func CreateTarContext(sourceDir string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, file); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}
