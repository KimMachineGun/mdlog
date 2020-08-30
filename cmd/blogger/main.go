package main

import (
	"log"
	"os"
)

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
