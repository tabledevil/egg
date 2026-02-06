package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	inputFile  = "questions.json"
	outputFile = "pkg/data/embedded.go"
	xorKey     = 0xAA // Simple XOR key
)

func main() {
	// Read JSON
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	// Obfuscate
	obfuscated := make([]byte, len(data))
	for i, b := range data {
		obfuscated[i] = b ^ xorKey
	}

	// Generate Go code
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, "package data")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "var embeddedData = []byte{")
	for i, b := range obfuscated {
		if i%12 == 0 {
			fmt.Fprint(f, "\n\t")
		}
		fmt.Fprintf(f, "0x%02x, ", b)
	}
	fmt.Fprintln(f, "\n}")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "func LoadRawData() []byte {")
	fmt.Fprintln(f, "\tdecoded := make([]byte, len(embeddedData))")
	fmt.Fprintln(f, "\tfor i, b := range embeddedData {")
	fmt.Fprintf(f, "\t\tdecoded[i] = b ^ 0x%02x\n", xorKey)
	fmt.Fprintln(f, "\t}")
	fmt.Fprintln(f, "\treturn decoded")
	fmt.Fprintln(f, "}")

	fmt.Printf("Successfully packed %d bytes into %s\n", len(data), outputFile)
}
