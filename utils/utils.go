// Contains function to be uesd across the application.
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var logPanic = log.Panic

func HandleErr(err error) {
	if err != nil {
		logPanic(err)
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

// 바이트에서 변수로 저장
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

func ToJSON[T any](payload T) []byte {
	json, err := json.Marshal(payload)
	HandleErr(err)
	return json
}

func GetPort() string {
	port := os.Args[2][6:]
	return port
}