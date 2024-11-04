package main

import (
	"encoding/json"
	"fmt"
	"github.com/jackpal/bencode-go"
	"os"
	"strings"
)

func main() {
	command := os.Args[1]
	// fmt.Println("Command:", command)

	switch command {
	case "decode":
		// fmt.Println("Decoding...")
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Printf("Error (%v), decoding given input: %s", err, bencodedValue)
			return
		}

		jsonOutput, err := json.Marshal(decoded)
		if err != nil {
			fmt.Printf("Error (%v), converting decoded data to JSON: %v", err, decoded)
			return
		}
		fmt.Println(string(jsonOutput))
	}
}

func decodeBencode(bencodedString string) (interface{}, error) {
	bencodedReader := strings.NewReader(bencodedString)
	res, err := bencode.Decode(bencodedReader)
	if err != nil {
		return "", err
	}

	return res, nil
}
