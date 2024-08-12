package main

import (
	"fmt"
	"os"

	"splitwise/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the input file path")
		return
	}

	filePath := os.Args[1]

	if err := cmd.ProcessFile(filePath); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
	}
}
