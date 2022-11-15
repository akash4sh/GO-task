package main

import (
	"io"
	"log"
	"os"
)

func Copy(src, dest string) {

	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	numBytes, err := io.Copy(destFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully copied %d bytes", numBytes)
	err = destFile.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Copy("sample.txt", "test.txt")
}
