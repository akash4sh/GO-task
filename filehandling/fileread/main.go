package main

import (
	"log"
	"os"
)

func FileRead(filePath string) {
	oFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer oFile.Close()
	buff := make([]byte, 100)
	for no, err := oFile.Read(buff); err == nil; no, err = oFile.Read(buff) {
		if no > 0 {
			os.Stdout.Write(buff[0:no])
		}
	}
}

func main() {
	FileRead("sample.txt")
}
