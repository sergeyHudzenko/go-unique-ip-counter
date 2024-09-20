package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

// Sort a chunk of IP addresses concurrently and write to a temporary file
func sortAndWriteChunk(chunk []string, chunkNum int, wg *sync.WaitGroup, tempFiles chan<- string) {
	defer wg.Done()
	sort.Strings(chunk)
	fileName := fmt.Sprintf("temp_chunk_%d.txt", chunkNum)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, ip := range chunk {
		_, err := writer.WriteString(ip + "\n")
		if err != nil {
			fmt.Println("Error writing to temp file:", err)
			return
		}
	}
	writer.Flush()

	// Send the temp file name to the channel
	tempFiles <- fileName
}

// MinHeap to merge sorted chunks
type FileIP struct {
	ip     string
	reader *bufio.Reader
	index  int
}

type MinHeap []*FileIP

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].ip < h[j].ip }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(*FileIP))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// Merge sorted chunks and count unique IP addresses
func mergeChunks(tempFiles []string) (int, error) {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	// Open all chunk files
	for i, fileName := range tempFiles {
		file, err := os.Open(fileName)
		if err != nil {
			return 0, err
		}
		reader := bufio.NewReader(file)
		ip, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, err
		}
		ip = strings.TrimSpace(ip)
		if len(ip) > 0 {
			heap.Push(minHeap, &FileIP{ip: ip, reader: reader, index: i})
		}
	}

	uniqueCount := 0
	var lastIP string

	for minHeap.Len() > 0 {
		minIP := heap.Pop(minHeap).(*FileIP)

		// Count only unique IPs
		if minIP.ip != lastIP {
			uniqueCount++
			lastIP = minIP.ip
		}

		// Read the next IP from the same file
		nextIP, err := minIP.reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, err
		}
		nextIP = strings.TrimSpace(nextIP)
		if len(nextIP) > 0 {
			minIP.ip = nextIP
			heap.Push(minHeap, minIP)
		}
	}

	return uniqueCount, nil
}

func main() {
	const chunkSize = 1000000 // Adjust based on memory capacity

	fileName := "large_ip_file.txt"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	tempFilesChan := make(chan string, 10) // Buffered channel to hold temp file names

	chunk := []string{}
	chunkNum := 0
	tempFiles := []string{}

	// Start a goroutine to collect temp file names
	go func() {
		for tempFile := range tempFilesChan {
			tempFiles = append(tempFiles, tempFile)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) >= chunkSize {
			wg.Add(1)
			go sortAndWriteChunk(chunk, chunkNum, &wg, tempFilesChan)
			chunk = []string{}
			chunkNum++
		}
	}

	// Process the last chunk if it exists
	if len(chunk) > 0 {
		wg.Add(1)
		go sortAndWriteChunk(chunk, chunkNum, &wg, tempFilesChan)
	}

	wg.Wait()              // Wait for all sorting goroutines to finish
	close(tempFilesChan)   // Close the channel to signal the end of temp files

	// Merge sorted chunks and count unique IP addresses
	uniqueCount, err := mergeChunks(tempFiles)
	if err != nil {
		fmt.Println("Error merging chunks:", err)
		return
	}

	fmt.Printf("Number of unique IP addresses: %d\n", uniqueCount)

	// Clean up temporary files
	for _, fileName := range tempFiles {
		os.Remove(fileName)
	}
}
