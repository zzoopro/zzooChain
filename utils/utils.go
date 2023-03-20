package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Hash[T any](data T) string {
	dataAsString := fmt.Sprintf("%v", data)
	hash := sha256.Sum256([]byte(dataAsString))
	return fmt.Sprintf("%x",hash)
}

func ToBytes[T any](i T) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes[T any](target T, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(target))
}

func Splitter(s ,sep string, i int) string {
	result := strings.Split(s, sep)
	if len(result) - 1 < i {
		return ""
	}
	return result[i]
}