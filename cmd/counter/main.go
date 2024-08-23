package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	bitArraySize = 2000000000 // Bit array size for the Blum filter (the bigger, the more accurate)
	hashFuncs    = 5          // Number of hash functions for the Blum filter (the bigger, the more accurate)
)

type BloomFilter struct {
	bitArray []uint64
}

// NewBloomFilter - creating a new Blum filter
func NewBloomFilter(size int) *BloomFilter {
	arraySize := (size + 63) / 64
	return &BloomFilter{bitArray: make([]uint64, arraySize)}
}

func hash1(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}

func hash2(data string) uint32 {
	h := fnv.New32()
	h.Write([]byte(data))
	return h.Sum32()
}

func hash3(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32() + uint32(len(data))
}

// Add - adding IP to the Blum filter
func (bf *BloomFilter) Add(ip string) {
	indices := bf.getIndices(ip)
	for _, idx := range indices {
		wordIndex := idx / 64
		bitIndex := idx % 64
		atomic.OrUint64(&bf.bitArray[wordIndex], 1<<bitIndex)
	}
}

// Test - checking for IP in the Blum filter
func (bf *BloomFilter) Test(ip string) bool {
	indices := bf.getIndices(ip)
	for _, idx := range indices {
		wordIndex := idx / 64
		bitIndex := idx % 64
		if atomic.LoadUint64(&bf.bitArray[wordIndex])&(1<<bitIndex) == 0 {
			return false
		}
	}

	return true
}

// getIndices - getting indices for the Blum filter
func (bf *BloomFilter) getIndices(ip string) []int {
	indices := make([]int, hashFuncs)
	h1 := hash1(ip)
	h2 := hash2(ip)
	h3 := hash3(ip)
	indices[0] = int(h1) % (len(bf.bitArray) * 64)
	indices[1] = int(h2) % (len(bf.bitArray) * 64)
	indices[2] = int(h3) % (len(bf.bitArray) * 64)

	return indices
}

// workerBloomFilter - worker for the Blum filter algorithm
func workerBloomFilter(id int, bf *BloomFilter, jobs <-chan string, wg *sync.WaitGroup, result *int32) {
	defer wg.Done()
	for ip := range jobs {
		if !bf.Test(ip) {
			bf.Add(ip)
			atomic.AddInt32(result, 1)
		}
	}
}

// workerSimple - worker for the simple "naive" algorithm
func workerSimple(id int, uniqueIPs *sync.Map, jobs <-chan string, wg *sync.WaitGroup, result *int32) {
	defer wg.Done()
	for ip := range jobs {
		if _, loaded := uniqueIPs.LoadOrStore(ip, struct{}{}); !loaded {
			atomic.AddInt32(result, 1)
		}
	}
}

// parallelReadFile - reading file in parallel
func parallelReadFile(filePath string, jobs chan<- string, numWorkers int) {
	defer close(jobs)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error when opening a file:", err)

		os.Exit(1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)

		os.Exit(1)
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			offset := int64(workerId) * chunkSize
			end := offset + chunkSize

			if workerId == numWorkers-1 {
				end = fileSize
			}

			buffer := make([]byte, end-offset)
			file.ReadAt(buffer, offset)

			scanner := bufio.NewScanner(strings.NewReader(string(buffer)))

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					jobs <- line
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("File read error:", err)
			}
		}(i)
	}

	wg.Wait()
}

func main() {
	start := time.Now()

	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d workers\n", numWorkers)

	inputFile := "result.txt"
	algorithm := 1

	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		var err error
		algorithm, err = strconv.Atoi(os.Args[2])
		if err != nil || (algorithm != 1 && algorithm != 2) {
			fmt.Println("Incorrect algorithm. Use `1` for the Blum filter or `2` for the simple algorithm.")

			os.Exit(1)
		}
	}

	jobs := make(chan string, 100)
	var wg sync.WaitGroup
	var uniqueCount int32

	if algorithm == 1 {
		fmt.Println("Using the Blum filter")
		bf := NewBloomFilter(bitArraySize)
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go workerBloomFilter(w, bf, jobs, &wg, &uniqueCount)
		}
	} else if algorithm == 2 {
		fmt.Println("Using a simple `naive` algorithm")
		uniqueIPs := &sync.Map{}
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go workerSimple(w, uniqueIPs, jobs, &wg, &uniqueCount)
		}
	}

	parallelReadFile(inputFile, jobs, numWorkers)

	wg.Wait()

	fmt.Printf("Total number of unique IP addresses: %d\n", uniqueCount)
	fmt.Printf("Execution time: %v\n", time.Since(start))
}
