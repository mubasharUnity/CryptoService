package main

import (
	"encoding/binary"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
}

func (bch *BcryptHasher) Process(password []byte) ([]byte, error) {
	return bch.Hasher(password)
}

func (bch *BcryptHasher) Hasher(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

type BcryptHashComparer struct {
}

func (bch *BcryptHashComparer) Process(payload []byte) ([]byte, error) {
	passwordLength := binary.BigEndian.Uint16(payload[0:2])
	passwordEndByteIndex := 2 + passwordLength
	password := payload[2:passwordEndByteIndex]

	hashLength := binary.BigEndian.Uint16(payload[passwordEndByteIndex : passwordEndByteIndex+2])
	hashStartByteIndex := passwordEndByteIndex + 2
	hash := payload[hashStartByteIndex : hashStartByteIndex+hashLength]
	return bch.HashComparer(hash, password)
}

func (*BcryptHashComparer) HashComparer(hashedPassword, password []byte) ([]byte, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	var isSucces uint8 = 0x00
	if err == nil {
		isSucces = 0x01
	}

	return []byte{isSucces}, err
}
