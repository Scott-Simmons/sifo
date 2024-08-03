package decrypt

// Ref: https://medium.com/insiderengineering/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes-ceb598a34f41

import (
	"crypto/aes"
  "SecureSyncDrive/pkg/encrypt"
	"crypto/cipher"
	"fmt"
	"os"
)


func Decrypt(filePathToDecrypt string, privateKeyPath string) (error) {

  encryptedData, err := os.ReadFile(filePathToDecrypt)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

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
  // TODO: This might be the source of the bug...
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
}
