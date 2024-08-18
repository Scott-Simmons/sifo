package test_helpers

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"testing"
)

func CreateTarFile(t *testing.T, tarFileNameToWriteTo string, contentToWrite string, txtFileToWriteTo string) error {
	tarFile, err := os.CreateTemp("", tarFileNameToWriteTo)
	if err != nil {
		return err
	}
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()
	header := &tar.Header{
		Name:     txtFileToWriteTo,
		Size:     int64(len(contentToWrite)),
		Mode:     0600,
		Typeflag: tar.TypeReg,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tarWriter.Write([]byte(contentToWrite)); err != nil {
		return err
	}
	return nil
}

func CreateGzipTarFile(t *testing.T, tarFileNameToWriteTo string, contentToWrite string, txtFileToWriteTo string) error {
	tarFile, err := os.CreateTemp("", tarFileNameToWriteTo)
	if err != nil {
		return err
	}
	gzipWriter := gzip.NewWriter(tarFile)
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	defer gzipWriter.Close()
	header := &tar.Header{
		Name:     txtFileToWriteTo,
		Size:     int64(len(contentToWrite)),
		Mode:     0600,
		Typeflag: tar.TypeReg,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tarWriter.Write([]byte(contentToWrite)); err != nil {
		return err
	}
	return nil
}
