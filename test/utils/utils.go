package utils

import (
	"fmt"
	"io/ioutil"
	"log"
)

// GetFilesInPwd lists files in the pwd
func GetFilesInPwd() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
