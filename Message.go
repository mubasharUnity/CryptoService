package main

import (
	"encoding/binary"
	"errors"
)

type MessageProcessor interface {
	Process(data []byte) ([]byte, error)
}
type MessageHeader struct {
	msgVersion    uint8
	mode          uint8
	msgId         string
	msgIdBytes    []byte
	lengthPayload uint16
}
type Message struct {
	MessageHeader
	payload []byte
}
type OutMessage struct {
	inMessage       *Message
	processedBuffer []byte
}

func (msgInstanceHeader *MessageHeader) ValidateHeader() bool {
	if msgInstanceHeader.msgVersion != MSG_VERSION {
		println("Msg version not supported.")
		return false
	}
	if (msgInstanceHeader.mode & 0x10) != 0 {
		println("Requested Mode Invalid, expected request mode")
		return false
	}
	if len(msgInstanceHeader.msgId) == 0 {
		println("Message ID missing")
		return false
	}
	return true
}
func DeconstructMessage(buffer []byte, readByteLength int) (msg *Message, err error) {
	readStatus := msg.CheckIfReadOK(readByteLength)
	if !readStatus {
		return nil, errors.New("msg is either too short or too long")
	}

	msg = new(Message)
	msg.msgVersion = uint8(buffer[0])

	msg.mode = uint8(buffer[1])

	msg.msgIdBytes = buffer[2:18]
	msg.msgId = string(msg.msgIdBytes)

	headerValidated := msg.ValidateHeader()
	if !headerValidated {
		return nil, errors.New("incorrect header")
	}

	msg.lengthPayload = binary.BigEndian.Uint16(buffer[18:20])

	if !msg.ValidatePayload(readByteLength) {
		return nil, errors.New("invalid payload length")
	}
	msg.payload = buffer[20 : 20+msg.lengthPayload]

	return
}
func (msgInstance *Message) ValidatePayload(readBufferLength int) bool {
	diff := readBufferLength - HEADER_LENGTH - int(msgInstance.lengthPayload)
	return diff == 0
}
func (msgInstance *Message) CheckIfReadOK(n int) bool {
	if n < HEADER_LENGTH {
		println("Msg has wrong header")
		return false
	}
	if n == BUFFER_LENGTH {
		println("Msg too long(server buffer overload).")
		return false
	}
	return true
}

func (msg *OutMessage) FormResponseBuffer() []byte {
	processedBufferLength := len(msg.processedBuffer)
	if processedBufferLength == 0 {
		return nil
	}

	returnBuffer := make([]byte, 0, HEADER_LENGTH+processedBufferLength)
	modeBytes := []byte{msg.inMessage.mode | 0x80}
	returnPayloadLength := make([]byte, 2)
	binary.BigEndian.PutUint16(returnPayloadLength, uint16(len(msg.processedBuffer)))
	returnBuffer = append(returnBuffer, msg.inMessage.msgVersion)
	returnBuffer = append(returnBuffer, modeBytes...)
	returnBuffer = append(returnBuffer, msg.inMessage.msgIdBytes...)
	returnBuffer = append(returnBuffer, returnPayloadLength...)
	returnBuffer = append(returnBuffer, msg.processedBuffer...)

	return returnBuffer
}
