package main

import "sync"

var (
	factoryOnce     sync.Once
	processorsSlice []MessageProcessor
)

func configureProcessors() {
	aesEncyptor := AESEncyptor{}
	aesDecryptor := AESDecyptor{}
	bcryptHasher := BcryptHasher{}
	bcryptHashComparer := BcryptHashComparer{}

	processorsSlice = []MessageProcessor{&aesEncyptor, &aesDecryptor, &bcryptHasher, &bcryptHashComparer}
}

func GetProcessorForMode(mode uint8) MessageProcessor {
	factoryOnce.Do(configureProcessors)
	return processorsSlice[mode]
}
