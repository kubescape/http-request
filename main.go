package main

import (
	"fmt"
	"log"
)

func main() {
	parser := NewFlagParser()
	parser.parser()
	if err := parser.validate(); err != nil {
		log.Fatal(err)
	}

	if response, err := Request(parser); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("response: %s", response)
	}

}
