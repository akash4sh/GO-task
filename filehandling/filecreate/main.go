package main

import (
	"log"
	"os"
)

func CreateEmptyFile() {
	myFile, err := os.Create("sample.txt")
	if err != nil {
		log.Fatal("ERROR! ", err)
	}
	log.Println("Empty file created successfully. ", myFile)
	myFile.Close()
}

func main() {
	CreateEmptyFile()
}
