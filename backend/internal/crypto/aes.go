package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const fixedIV = "wjp5i4U9bZZAJLGs"

func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// PKCS#7 padding
	padLen := aes.BlockSize - len(plaintext)%aes.BlockSize
	padded := append(plaintext, bytes.Repeat([]byte{byte(padLen)}, padLen)...)

	iv := []byte(fixedIV)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return ciphertext, nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	iv := []byte(fixedIV)
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS#7 padding
	padLen := int(plaintext[len(plaintext)-1])
	if padLen <= 0 || padLen > aes.BlockSize {
		return nil, errors.New("invalid padding size")
	}
	return plaintext[:len(plaintext)-padLen], nil
}
