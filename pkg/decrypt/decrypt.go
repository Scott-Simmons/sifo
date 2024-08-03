package decrypt

// Ref: https://medium.com/insiderengineering/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes-ceb598a34f41

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/aes"
  "SecureSyncDrive/pkg/encrypt"
	"crypto/cipher"
  "path/filepath"
	"fmt"
	"io"
	"os"
)


func Decrypt(filePathToDecrypt string, privateKeyPath string) (error) {

  encryptedData, err := os.ReadFile(filePathToDecrypt)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

  // TODO: Read key from file
  privateKey, err := encrypt.ReadAES256KeyFromFile(privateKeyPath)
  if err != nil {
    fmt.Println("Failed to read key:", err)
    return err
  }

  block, err := aes.NewCipher(privateKey)
  if err != nil {
		return fmt.Errorf("failed to create AES cipher: %v", err)
	}
  
  // the IV was put at start of the ciphertext
  iv := encryptedData[:aes.BlockSize]
  ciphertext := encryptedData[aes.BlockSize:]

  mode := cipher.NewCBCDecrypter(block, iv)
  decrypted := make([]byte, len(ciphertext))
  mode.CryptBlocks(decrypted, ciphertext)

  // Unpad
  fmt.Println("Unpadding")
  unpadded, err := encrypt.Unpad(decrypted, aes.BlockSize)
  if err != nil {
    return err
  }

  outputFile := filePathToDecrypt + ".dec"
	err = os.WriteFile(outputFile, unpadded, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data to file: %v", err)
	}
  fmt.Printf("Decrypted data to file: %v", outputFile)

  return nil



//  fmt.Println("Extracting tar archive")
  // TODO: Should probably have this in a separate function
  // Extract tar
//  err = extractTar(unpadded)
//  if err != nil {
//    return err
 // }

  return nil

}


// TODO: Maybe move to somewhere else
func extractTar(tarData []byte) error {
  // assumes gzipped
  gzipReader, err := gzip.NewReader(bytes.NewReader(tarData))
  if err != nil {
    return err
  }
  defer gzipReader.Close()

  tarReader := tar.NewReader(gzipReader)

  for {
    header, err := tarReader.Next()
    if err == io.EOF {
      break
    }
    if err != nil {
      return err
    }

    switch header.Typeflag {
    case tar.TypeDir:
      if err := os.MkdirAll(header.Name, os.FileMode(header.Mode)); err != nil {
          return err
        }
    case tar.TypeReg:
      if err := os.MkdirAll(filepath.Dir(header.Name), 0755); err != nil {
          return err
        }
      outFile, err := os.Create(header.Name) 
      if err != nil {
        return err
      }
      defer outFile.Close()
      if _, err := io.Copy(outFile, tarReader); err != nil {
        return err
      }
    default:
      return fmt.Errorf("Tar entry type not supported: %v", header.Typeflag)
    }
  }
    return nil
  }
