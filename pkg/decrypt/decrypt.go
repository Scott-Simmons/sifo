package decrypt

// Ref: https://medium.com/insiderengineering/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes-ceb598a34f41

import (
	"SecureSyncDrive/pkg/encrypt"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"os"
)

// TODO: Refactor out constants

// TODO: Better error messages
func DoDecryptData(ciphertext []byte, privateKey []byte, iv []byte) ([]byte, error) {
	const pKeyLen = 32
	if len(privateKey) != pKeyLen {
		return nil, fmt.Errorf("Key length too short")
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("IV length too short")
	}
	block, err := aes.NewCipher(privateKey)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext, err = PKCS7Unpadding(plaintext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func PKCS7Unpadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	padding := data[len(data)-1]
	if int(padding) > len(data) {
		return nil, errors.New("padding size is larger than data")
	}
	return data[:len(data)-int(padding)], nil
}

func DecryptData(encryptedData []byte, privateKey []byte) ([]byte, error) {
	if len(encryptedData) < aes.BlockSize {
		return nil, fmt.Errorf("File too short to contain IV and encrypted data")
	}
	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]
	return DoDecryptData(ciphertext, privateKey, iv)
}

func Decrypt(filePathToDecrypt string, privateKeyPath string) error {
	encryptedData, err := os.ReadFile(filePathToDecrypt)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	privateKey, err := encrypt.ReadAES256KeyFromFile(privateKeyPath)
	if err != nil {
		fmt.Println("Failed to read key:", err)
		return err
	}
	decryptedData, err := DecryptData(encryptedData, privateKey)
	outputFile := filePathToDecrypt + ".dec"
	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data to file: %v", err)
	}
	fmt.Printf("Decrypted data to file: %v", outputFile)
	return nil
}

// TODO: Need to refactor and clean
