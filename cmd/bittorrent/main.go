package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/jackpal/bencode-go"
	"os"
	"strings"
)

type TorrentFile struct {
	bencodeMap map[string]interface{}
	infoMap    map[string]interface{}
	encodedMap bytes.Buffer
}

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

	case "info":

		torrentFileInfo := extractTorrentFileInformation(2)

		// calculate SHA-1 hash of the bencoded info dictionary
		infoHash := sha1.Sum(torrentFileInfo.encodedMap.Bytes())

		pieceLength := torrentFileInfo.infoMap["piece length"].(int64)
		pieces := torrentFileInfo.infoMap["pieces"].(string)

		fmt.Printf("Tracker URL: %v\n", torrentFileInfo.bencodeMap["announce"].(string))
		fmt.Printf("Length: %v\n", torrentFileInfo.infoMap["length"].(int64))
		fmt.Printf("Info Hash: %x\n", infoHash)
		fmt.Printf("Piece Length: %v\n", pieceLength)
		fmt.Printf("Piece Hashes:\n")
		var hash string
		for i := 0; i < len(pieces); i += 20 {
			hash = fmt.Sprintf("%x", pieces[i:i+20])
			fmt.Printf("%v\n", hash)
		}
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

func extractTorrentFileInformation(index int) *TorrentFile {
	fileName := os.Args[index]
	// fmt.Println("Torrent file: ", fileName)
	// parse the input file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	// decode the bencoded data
	bencodeData, err := bencode.Decode(file)
	if err != nil {
		fmt.Println("Error decoding bencode data: ", err)
		os.Exit(1)
	}

	// convert bencoded data a map
	bencodeMap, ok := bencodeData.(map[string]interface{})
	if !ok {
		fmt.Println("Error converting bencode data to map")
		os.Exit(1)
	}

	// get the info map
	infoMap, ok := bencodeMap["info"].(map[string]interface{})
	if !ok {
		fmt.Println("Error converting info data to map")
		os.Exit(1)
	}

	// peerAddress := os.Args[3]
	// fmt.Println("Peer address: ", peerAddress)

	var encodedInfoMap bytes.Buffer
	err = bencode.Marshal(&encodedInfoMap, infoMap)
	if err != nil {
		fmt.Println("Error marshalling info data: ", err)
		os.Exit(1)
	}

	return &TorrentFile{
		bencodeMap: bencodeMap,
		infoMap:    infoMap,
		encodedMap: encodedInfoMap,
	}
}
