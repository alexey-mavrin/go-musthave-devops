package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexey-mavrin/go-musthave-devops/internal/crypt"
)

func usage() string {
	return fmt.Sprintf("Usage: %s key_dir\n", os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(usage())
	}
	dir := os.Args[1]
	if err := crypt.GenerateStoreKeys(dir); err != nil {
		log.Fatal(err)
	}
}
