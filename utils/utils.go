package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Hash(payloads... string) string {
	text := ""	
	for _,payload := range payloads {
		text = text + payload
	}	
	hash := fmt.Sprintf("%x",sha256.Sum256([]byte(text)))	
	return hash
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}