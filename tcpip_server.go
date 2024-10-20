package main

import (
	"context"
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
	outChan := make(chan *OutMessage, 16)
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		conn.Close()
		for value := range outChan {
			fmt.Println("drained buffer of id: ", len(value.inMessage.msgId))
		}

		close(outChan)
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	go SendResponses(conn, outChan, ctx)

	for {
		msg := make([]byte, BUFFER_LENGTH)
		n, err := conn.Read(msg)
		if err != nil {
			println(err.Error())
			return
		}
		println(fmt.Sprintf("%d %s", n, err))

		msgInstance, errparsing := DeconstructMessage(msg, n)
		if errparsing != nil {
			fmt.Println(errparsing)
			return
		}
		fmt.Printf("new msg from %s with ID %s with mode %x\n", conn.RemoteAddr().String(), msgInstance.msgIdBytes, msgInstance.mode)

		go ProcessInputMessage(msgInstance, outChan)
	}
}

func ProcessInputMessage(inMsg *Message, chOutWriteOnly chan<- *OutMessage) {
	processed, err := GetProcessorForMode(inMsg.mode).Process(inMsg.payload)
	if err != nil {
		println(err)
	}
	outMsg := OutMessage{inMessage: inMsg, processedBuffer: processed}
	chOutWriteOnly <- &outMsg
}

func SendResponses(conn net.Conn, outputChan <-chan *OutMessage, ctx context.Context) {
	for {
		select {
		case responseMsg := <-outputChan:
			responseBuffer := responseMsg.FormResponseBuffer()
			conn.Write(responseBuffer)
		case <-ctx.Done():
			fmt.Println("Context canceled.")
			return
		}
	}
}
