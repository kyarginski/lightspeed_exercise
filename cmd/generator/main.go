package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()

	defaultFileName := "result.txt"
	defaultSizeMB := 10

	var fileName string
	var sizeMB int

	if len(os.Args) > 1 {
		fileName = os.Args[1]
	} else {
		fileName = defaultFileName
	}

	// File size in MB (if defined)
	if len(os.Args) > 2 {
		var err error
		sizeMB, err = strconv.Atoi(os.Args[2])
		if err != nil || sizeMB <= 0 {
			fmt.Println("Incorrect file size value, default value is used:", defaultSizeMB, "MB")
			sizeMB = defaultSizeMB
		}
	} else {
		sizeMB = defaultSizeMB
	}

	// File size in bytes
	sizeBytes := sizeMB * 1024 * 1024

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error when creating a file:", err)

		os.Exit(1)
	}
	defer file.Close()

	seed, _ := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	rng := mathrand.New(mathrand.NewSource(seed.Int64()))

	totalBytes := 0

	for totalBytes < sizeBytes {
		// Generating a random IPv4 address
		ip := fmt.Sprintf(
			"%d.%d.%d.%d\n",
			rng.Intn(256), rng.Intn(256),
			rng.Intn(256), rng.Intn(256),
		)

		bytesWritten, err := file.WriteString(ip)
		if err != nil {
			fmt.Println("Error when writing to a file:", err)

			os.Exit(1)
		}

		totalBytes += bytesWritten
	}

	fmt.Printf("File `%s` was successfully created with size %d MB\n", fileName, sizeMB)

	fmt.Printf("Execution time: %v\n", time.Since(start))
}
