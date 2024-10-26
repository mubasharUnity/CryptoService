package main

import (
	"reflect"
	"testing"
)

func TestAESEncyptionCycle(t *testing.T) {
	enc := AESEncyptor{}
	dec := AESDecyptor{}
	inputData := []byte{0x11, 0xae, 0x43}

	encResult, _ := enc.Process(inputData)
	decResult, _ := dec.Process(encResult)

	if !reflect.DeepEqual(inputData, decResult) {
		t.Errorf("Encryption and Decryption do not match")
	}
}
