package encrypt

// Ref: https://earthly.dev/blog/cryptography-encryption-in-go/

import (
	"SecureSyncDrive/pkg/archive"
  "encoding/hex"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
//	"errors"
	"io"
	"os"
  "fmt"
)

func GenerateKey() ([]byte, error) {
  aes256ByteSize := 32
	key := make([]byte, aes256ByteSize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
  fmt.Printf("Generated key: %x\n", key)
	return key, nil
}

func trimNewLinesFromBytes(btyes []byte) ([]byte) {
  const newlineLF = 0x0A
  const newlineCR = 0x0D
  var cleanedData []byte
	for _, byte := range btyes {
		if byte != newlineLF && byte != newlineCR {
			cleanedData = append(cleanedData, byte)
		}
	}
  return cleanedData
}

// TODO: put this in another place
func ReadAES256KeyFromFile(keyFilePath string) ([]byte, error) {
	keyHex, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(string(trimNewLinesFromBytes(keyHex)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hexadecimal key: %w", err)
	}

  const expectedBtyes = 32

	if len(key) != expectedBtyes {
		return nil, fmt.Errorf("key length must be %d bytes for AES-256 but found %d bytes", expectedBtyes, (len(key)))
	}
	return key, nil
}


func EncryptTarBall(tarBallToEncrypt string, encryptedTarballPath string, privateKeyPath string) error {

    key, err := ReadAES256KeyFromFile(privateKeyPath)
    if err != nil {
      fmt.Println("Failed to read key:", err)
      return err
    }
  
    srcFile, err := os.Open(tarBallToEncrypt)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    // TODO: isTar has a bug in it... need to fix later for now just don't use it
    _, err = archive.FileIsTar(srcFile)
    if err != nil {
      return err
    }
    //if !isTar {
    //  return errors.New("file is not a valid tar archive")
    //}

    encryptedTarballFile, err := os.Create(encryptedTarballPath)
    if err != nil {
        return err
    }
    defer encryptedTarballFile.Close()

    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    // Create initialisation vector and fill it with random numbers
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return err
    }

    // Write the IV to the start of the destination file
    if _, err := encryptedTarballFile.Write(iv); err != nil {
        return err
    }

    // Create CBC mode cipher encrypter function
    encrypter := cipher.NewCBCEncrypter(block, iv)

    // Begin the encryption
    bufferSize := 1024 // arbitrary
    buf := make([]byte, bufferSize)
    for {
        // File content read into buffer of fixed length
        n, err := srcFile.Read(buf)
        if err != nil && err != io.EOF {
            return err
        }
        if n == 0 {
            break // EOF reached
        }

        // Pad data to be a multiple of the block size
        paddedData := pad(buf[:n], aes.BlockSize)
        
        // Encrypt and write to the destination file
        // Do the encryption in place, overwriting the padded data
        encrypter.CryptBlocks(paddedData, paddedData)
        if _, err := encryptedTarballFile.Write(paddedData); err != nil {
            return err
        }
    }

    return nil
}

func pad(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

// TODO: Put this in another place
func Unpad(data []byte, blockSize int) ([]byte, error) {
    length := len(data)
    if length == 0 {
        return nil, fmt.Errorf("data length is zero")
    }

    padding := data[length-1]
    if int(padding) > blockSize || int(padding) > length {
        return nil, fmt.Errorf("invalid padding size")
    }

    for _, b := range data[length-int(padding):] {
        if b != padding {
            return nil, fmt.Errorf("invalid padding byte")
        }
    }

    return data[:length-int(padding)], nil
}
