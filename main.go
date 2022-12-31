package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	srcFile := os.Args[1]
	dstDir := os.Args[2]
	newDir := ""
	if len(os.Args) > 3 {
		newDir = os.Args[3]
	}

	fr, err := os.Open(srcFile)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	if !fileExists(dstDir) {
		err := os.Mkdir(dstDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		dName := hdr.Name
		if newDir != "" {
			idx := strings.Index(hdr.Name, string(filepath.Separator))
			if idx != -1 {
				dName = strings.Replace(dName, dName[0:idx], newDir, 1)
			}
		}

		dstPath := filepath.Join(dstDir, dName)

		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(dstPath, os.FileMode(hdr.Mode)); err != nil {
				panic(err)
			}
			continue
		}

		if err := extractFile(tr, hdr, dstPath); err != nil {
			panic(err)
		}
	}
}

func extractFile(tr *tar.Reader, hdr *tar.Header, dstPath string) error {
	dstDir := filepath.Dir(dstPath)
	if !fileExists(dstDir) {
		err := os.MkdirAll(dstDir, 0755)
		if err != nil {
			return err
		}
	}

	fw, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed os.Create: %w", err)
	}
	defer fw.Close()

	if err := os.Chmod(dstPath, os.FileMode(hdr.Mode)); err != nil {
		return fmt.Errorf("failed to os.Chmod to %s %8d: %w", dstPath, hdr.Mode, err)
	}

	if _, err := io.Copy(fw, tr); err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
