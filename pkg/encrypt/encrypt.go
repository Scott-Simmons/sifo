package encrypt

// Ref: https://earthly.dev/blog/cryptography-encryption-in-go/
// Ref: https://gist.github.com/yingray/57fdc3264b1927ef0f984b533d63abab

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GenerateKey() (string, error) {
	aes256ByteSize := 32
	key := make([]byte, aes256ByteSize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	fmt.Printf("Generated key (hex): %x\n", key)
	return hex.EncodeToString(key), nil
}

func trimNewLinesFromBytes(btyes []byte) []byte {
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
	const expectedBtyes = 32
	keyHex, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(string(trimNewLinesFromBytes(keyHex)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hexadecimal key: %w", err)
	}
	if len(key) != expectedBtyes {
		return nil, fmt.Errorf("key length must be %d bytes for AES-256 but found %d bytes", expectedBtyes, (len(key)))
	}
	return key, nil
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func GetAES256Encrypted(plaintext []byte, iv []byte, privateKey []byte) ([]byte, error) {
	const pKeyLen = 32
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("IV length must be %d bytes, got %d", aes.BlockSize, len(iv))
	}
	if len(privateKey) != pKeyLen {
		return nil, fmt.Errorf("key length must be %d bytes, got %d", pKeyLen, len(privateKey))
	}
	block, err := aes.NewCipher(privateKey)
	if err != nil {
		return nil, err
	}
	plainTextBlock := PKCS7Padding(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainTextBlock)

	// Prepend IV to ciphertext
	result := append(iv, ciphertext...)

	return result, nil
}

func generateRandomIV(blockSize int) ([]byte, error) {
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	return iv, nil
}

func EncryptTarBall(tarBallToEncrypt string, encryptedTarballPath string, privateKeyPath string) error {
	key, err := ReadAES256KeyFromFile(privateKeyPath)
	if err != nil {
		fmt.Println("Failed to read key:", err)
		return err
	}
	plaintextBytes, err := os.ReadFile(tarBallToEncrypt)
	if err != nil {
		return fmt.Errorf("Failed to read file %s: %w", tarBallToEncrypt, err)
	}
	iv, err := generateRandomIV(aes.BlockSize)
	if err != nil {
		return fmt.Errorf("Failed to generate random IV: %w", err)
	}
	encryptedData, err := GetAES256Encrypted(plaintextBytes, iv, key)
	err = os.WriteFile(encryptedTarballPath, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write encrypted data to file: %w", err)
	}
	return nil
}
