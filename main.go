package main

import (
	"log"
)

func main() {
	parser := NewFlagParser()
	parser.parser()
	if err := parser.validate(); err != nil {
		log.Fatal(err)
	}

	if err := Request(parser); err != nil {
		log.Fatal(err)
	}

}
