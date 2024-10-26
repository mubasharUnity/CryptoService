package main

import (
	"reflect"
	"testing"
)

func TestAESEncyptionCycle(t *testing.T) {
	enc := AESEncyptor{}
	dec := AESDecyptor{}
	inputData := []byte{0x11, 0xae, 0x43}

	encResult, errE := enc.Process(inputData)
	if errE != nil {
		t.Error(errE.Error())
	}
	decResult, errD := dec.Process(encResult)
	if errD != nil {
		t.Error(errD.Error())
	}

	if !reflect.DeepEqual(inputData, decResult) {
		t.Error("Encryption and Decryption do not match")
	}
}
