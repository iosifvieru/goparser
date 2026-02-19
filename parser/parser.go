package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"
)

type LogEntry struct {
	ip       string
	respCode string
}

var logPattern = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[([^\]]+)\] "(\S+) ([^"]*) (\S+)" (\d{3}) (\S+) "([^"]*)" "([^"]*)" "([^"]*)"$`)

func mapLogChunkToLogEntries(lines []string) []LogEntry {
	logEntries := make([]LogEntry, 0, len(lines))

	for _, line := range lines {
		logEntry, err := parseLine(line)
		if err != nil {
			continue
		}
		logEntries = append(logEntries, logEntry)
	}

	return logEntries
}

func parseLine(line string) (LogEntry, error) {
	m := logPattern.FindStringSubmatch(line)
	if m == nil {
		return LogEntry{}, fmt.Errorf("line does not match log pattern: %q\n", line)
	}

	ipAddr := m[1]
	respCode := m[8]
	return LogEntry{ip: ipAddr, respCode: respCode}, nil
}

func readChunks(scanner *bufio.Scanner, chunkSize int, out chan<- []string) error {
	defer close(out)

	lines := make([]string, 0, chunkSize)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) >= chunkSize {
			chunk := make([]string, len(lines))
			copy(chunk, lines)

			out <- chunk
			lines = lines[:0]
		}
	}

	// send last partial chunk
	if len(lines) > 0 {
		chunk := make([]string, len(lines))
		copy(chunk, lines)
		out <- chunk
	}

	return scanner.Err()
}

func parseWorker(chunks <-chan []string, entries chan<- []LogEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	for chunk := range chunks {
		parsed := mapLogChunkToLogEntries(chunk)
		if len(parsed) > 0 {
			entries <- parsed
		}
	}
}

func GenerateReport(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("%q no such file or directory", filepath)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	chunkCh := make(chan []string, 16)
	entryCh := make(chan []LogEntry, 16)

	const N = 4096
	startTime := time.Now()

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	wg.Add(numWorkers)
	for range numWorkers {
		go parseWorker(chunkCh, entryCh, &wg)
	}

	go func() {
		wg.Wait()
		close(entryCh)
	}()

	producerErrCh := make(chan error, 1)
	go func() {
		producerErrCh <- readChunks(scanner, N, chunkCh)
	}()

	totalEntries := 0
	ipMap := make(map[string]int)
	statusCodesMap := make(map[string]int)
	for parsedBatch := range entryCh {
		totalEntries += len(parsedBatch)

		for _, entry := range parsedBatch {
			ipMap[entry.ip]++
			statusCodesMap[entry.respCode]++
		}
	}

	if err := <-producerErrCh; err != nil {
		return err
	}

	elapsed := time.Since(startTime)
	fmt.Println("Log source:", filepath)
	fmt.Println("Log parsing time:", elapsed)
	fmt.Println("Total requests:", totalEntries)

	for statusCode, count := range statusCodesMap {
		percentage := float64(count) / float64(totalEntries) * 100
		fmt.Fprintf(os.Stdout, "HTTP %s: %d requests (%f %%)\n", statusCode, count, percentage)
	}

	fmt.Println("Unique addresses:", len(ipMap))

	return nil
}
