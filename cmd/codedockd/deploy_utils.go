package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractArchiveTo(dest string, src io.Reader) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		tr := tar.NewReader(src)
		return untarTo(dest, tr)
	}
	defer gzr.Close()
	return untarTo(dest, tar.NewReader(gzr))
}

func untarTo(dest string, tr *tar.Reader) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		path := filepath.Join(dest, filepath.Clean(hdr.Name))
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(path, 0o755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(path), 0o755)
			f, _ := os.Create(path)
			if f != nil {
				_, _ = io.Copy(f, tr)
				f.Close()
			}
		}
	}
	return nil
}

func findSourceDir(dir string) string {
	entries, _ := os.ReadDir(dir)
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(dir, entries[0].Name())
	}
	return dir
}

func extractRepoName(url string) string {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "app"
}
