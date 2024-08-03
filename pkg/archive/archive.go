package archive

// Reference: https://www.arthurkoziel.com/writing-tar-gz-files-in-go/

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func FileIsTar(file *os.File) (bool, error) {
	const tarHeaderSizeBytes = 512
	headerBuffer := make([]byte, tarHeaderSizeBytes)

	_, err := file.Read(headerBuffer)
	if err != nil && err != io.EOF {
		return false, nil // don't throw error becuase it just means invalid tar
	}
	if err != nil {
		return false, nil
	}

	// Seek back to the start of the file to read it fully later
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("WTF")
		return false, err
	}

	tarReader := tar.NewReader(bytes.NewReader(headerBuffer))

	_, err = tarReader.Next()
	if err == io.EOF {
		// Empty file could be a true tar
		return true, nil
	}
	if err != nil {
		return false, nil // dont throw error because it just means invalid tar
	}

	// If header read then its a tar
	return true, nil

}

func ArchiveFolder(srcDirPath string, destTarPath string) error {
	destTar, err := os.Create(destTarPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", destTarPath, err.Error()))
	}
	defer destTar.Close()

	gzipWriter := gzip.NewWriter(destTar)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Iterate through the directory tree and add everything to the tar archive
	if err := addDirToTar(tarWriter, srcDirPath, srcDirPath); err != nil {
		return fmt.Errorf("Could not add write directory to tar: %w", err)
	}
	return nil
}

func addFileToTar(tarWriter *tar.Writer, filePath string, headerName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}
	// Need to overwrite header.Name because need complete path relative to the base directory to keep the tree structure... hence need to pass in the full path of the file and base dir. Since default fileinfo header name is just the file name not the full path.
	header.Name = headerName

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// Copying this data into the tar writer which is chained to the gzipwriter for compression
	_, err = io.Copy(tarWriter, file)
	return err
}

func addDirToTar(tarWriter *tar.Writer, baseDir string, dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	// Get file information
	fileInfo, err := d.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}

	// Validate if it is a directory
	if !fileInfo.IsDir() {
		fmt.Println("Not a directory")
		return err
	}

	// Read directory entries
	files, err := d.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Main logic - get relative path and full path
		fullPath := path.Join(dir, file.Name())
		relPath, _ := filepath.Rel(baseDir, fullPath)

		if file.IsDir() {
			if err := addDirToTar(tarWriter, baseDir, fullPath); err != nil {
				return err
			}
		} else {
			if err := addFileToTar(tarWriter, fullPath, relPath); err != nil {
				return err
			}
		}
	}

	return nil
}
