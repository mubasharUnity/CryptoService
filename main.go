package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
)

func getKeyPairPaths() (*string, *string) {
	key := os.Getenv("tls_key")
	cert := os.Getenv("tls_cert")
	if key == "" || cert == "" {
		return nil, nil
	}
	return &key, &cert
}

func main() {
	keyPath, certPath := getKeyPairPaths()
	cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	defer listener.Close()
	fmt.Println("Server on port 8443...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		go HandleConnection(conn)
	}
}
