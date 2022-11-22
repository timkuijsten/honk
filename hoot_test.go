package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestHooterize(t *testing.T) {
	dlog = log.Default()
	fd, err := os.Open("lasthoot.html")
	if err != nil {
		return
	}
	seen := make(map[string]bool)
	hoots := hootextractor(fd, "lasthoot.html", seen)
	fmt.Printf("hoots: %s\n", hoots)
}
