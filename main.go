package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/cedricmar/goover/page"
)

func usage() {
	fmt.Printf("Usage:\n  goover <url> <folder>\n")
	os.Exit(0)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	domain := os.Args[1]
	dest := os.Args[2]

	// @TODO - sanity check for url and folder
	// @TODO - create folder if not exist

	fmt.Printf("Hoovering %s in %s\n", domain, dest)

	var wg sync.WaitGroup
	var mux = sync.Mutex{}

	wg.Add(1)
	// @TODO - implement worker pool for big websites maybe ?
	go page.Get(domain, "/", dest, &wg, &mux)

	wg.Wait()

	fmt.Println("Done")
}
