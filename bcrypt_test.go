package main

import (
	"encoding/binary"
	"testing"
)

func TestHashingTestPossitive(t *testing.T) {
	hasher := BcryptHasher{}
	hashCompare := BcryptHashComparer{}
	password := []byte{0x31, 0x32, 0x33, 0x34}
	wrongPassword := []byte{0x41, 0x32, 0x33, 0x34}

	hash, _ := hasher.Process(password)

	hashLength := len(hash)
	hashLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(hashLengthBytes, uint16(hashLength))

	passwordLength := len(password)
	passwordLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(passwordLengthBytes, uint16(passwordLength))

	comparisonPayload := make([]byte, 0, hashLength+2+passwordLength+2)
	comparisonPayload = append(comparisonPayload, passwordLengthBytes...)
	comparisonPayload = append(comparisonPayload, password...)
	comparisonPayload = append(comparisonPayload, hashLengthBytes...)
	comparisonPayload = append(comparisonPayload, hash...)

	compareResults, _ := hashCompare.Process(comparisonPayload)

	if compareResults[0] == 0 {
		t.Error("Hash comparison with correct password failed")
	}

	wrongPasswordLength := len(wrongPassword)
	wrongPasswordLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(wrongPasswordLengthBytes, uint16(wrongPasswordLength))

	comparisonPayload = make([]byte, 0, hashLength+2+wrongPasswordLength+2)
	comparisonPayload = append(comparisonPayload, wrongPasswordLengthBytes...)
	comparisonPayload = append(comparisonPayload, wrongPassword...)
	comparisonPayload = append(comparisonPayload, hashLengthBytes...)
	comparisonPayload = append(comparisonPayload, hash...)

	compareResults, _ = hashCompare.Process(comparisonPayload)

	if compareResults[0] == 1 {
		t.Error("Hash comparison with wrong password passed")
	}
}
