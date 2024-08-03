package test_helpers

import (
  "compress/gzip"
  "testing"
  "os"
  "archive/tar"
  "io"
)


func CreateTarFile(t *testing.T, tarFileNameToWriteTo string, contentToWrite string, txtFileToWriteTo string) *os.File {
    tarFile, err := os.CreateTemp("", tarFileNameToWriteTo)
    if err != nil {
        t.Fatalf("Failed to create temp tar file: %v", err)
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
        t.Fatalf("Failed to write tar header: %v", err)
    }
    if _, err := tarWriter.Write([]byte(contentToWrite)); err != nil {
        t.Fatalf("Failed to write tar content: %v", err)
    }

    tarWriter.Close()
    tarFile.Seek(0, io.SeekStart)
    return tarFile
}

func CreateGzipTarFile(t *testing.T, tarFileNameToWriteTo string, contentToWrite string, txtFileToWriteTo string) *os.File {
    tarFile, err := os.CreateTemp("", tarFileNameToWriteTo)
    if err != nil {
        t.Fatalf("Failed to create temp tar file: %v", err)
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
        t.Fatalf("Failed to write tar header: %v", err)
    }
    if _, err := tarWriter.Write([]byte(contentToWrite)); err != nil {
        t.Fatalf("Failed to write tar content: %v", err)
    }

    // Close gzip writer to flush the data
    gzipWriter.Close()

    // Seek to the beginning of the file for reading
    tarFile.Seek(0, io.SeekStart)
    return tarFile
}
