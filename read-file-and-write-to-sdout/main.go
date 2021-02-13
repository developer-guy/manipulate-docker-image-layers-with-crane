package main

import (
	"io/ioutil"
	"log"
)

func main() {
	content, err := ioutil.ReadFile("hello-world.txt")

	if err != nil {
		log.Fatal("could not read file, error:", err)
	}

	log.Println("Content of the file is : ", string(content))
}
