package main

import (
	"log"

	"github.com/SociumR/bmlabs-test/internal/app/apiserver"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
