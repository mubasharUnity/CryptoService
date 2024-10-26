package main

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestMessageDecode(t *testing.T) {
	messageByte, err := hex.DecodeString("0100aabbaabbaabbaabbaabbaabbaabbaabb0003111111")
	if err != nil {
		t.Errorf(err.Error())
	}

	messageByteCount := len(messageByte)
	msgInstance, errparsing := DeconstructMessage(messageByte, messageByteCount)
	if errparsing != nil {
		t.Errorf(errparsing.Error())
	}

	if msgInstance.msgVersion != 0x01 {
		t.Errorf("Message version parsing failed")
	}
	if msgInstance.mode != 0x00 {
		t.Errorf("Message mode parsing failed")
	}
	if reflect.DeepEqual(msgInstance.msgIdBytes, messageByte[4:20]) {
		t.Errorf("Message id parsing failed")
	}
	if msgInstance.lengthPayload != 3 {
		t.Errorf("Payload length parsing failed")
	}

	if !msgInstance.ValidatePayload(messageByteCount) {
		t.Errorf("Payload length validation failed")
	}

	testPayload := messageByte[20:23]
	if !reflect.DeepEqual(msgInstance.payload, testPayload) {
		t.Errorf("Payload parsing failed")
	}

}
