package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ExportToFile(fileContent string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File saved: %s\n", filename)
}

func ImportFromFile(filename string) []byte{
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func FileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
