package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define flags
	name := flag.String("name", "World", "a name to greet")
	flag.Parse()

	// Print greeting
	fmt.Printf("Hello, %s!\n", *name)
}
