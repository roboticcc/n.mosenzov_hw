package main

import (
	"log"
	"os"
)

func fClose(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalf("unable to close %v", f.Name())
	}
}
