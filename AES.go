package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	once    sync.Once
	AES_KEY []byte
)

func initAESKey() []byte {
	once.Do(func() {
		key := os.Getenv("aes_key")
		AES_KEY = []byte(key)
	})
	return AES_KEY
}

type AESEncyptor struct {
}

func (enc *AESEncyptor) Process(data []byte) ([]byte, error) {
	return enc.encryptAES(data)
}

func (*AESEncyptor) encryptAES(payload []byte) ([]byte, error) {
	key := initAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, payload, nil)
	return cipherText, nil
}

type AESDecyptor struct {
}

func (dec *AESDecyptor) Process(data []byte) ([]byte, error) {
	return dec.decryptAES(data)
}

func (*AESDecyptor) decryptAES(cipherText []byte) ([]byte, error) {
	key := initAESKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
