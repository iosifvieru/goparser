package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

type LogEntry struct {
	ip       string
	respCode string
}

var logPattern = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[([^\]]+)\] "(\S+) ([^"]*) (\S+)" (\d{3}) (\S+) "([^"]*)" "([^"]*)" "([^"]*)"$`)

func parseLine(line string) (LogEntry, error) {
	m := logPattern.FindStringSubmatch(line)

	if m == nil {
		return LogEntry{}, fmt.Errorf("line does not match log pattern: %q\n", string(line))
	}

	// index 1 is ipv4 address
	ipAddr := m[1]

	// index 8 is HTTP status code in regex
	respCode := m[8]

	logEntry := LogEntry{ipAddr, respCode}
	return logEntry, nil
}

func printRequestReport(data map[string]int) {
	if len(data) == 0 {
		fmt.Println("the length of data is 0")
		return
	}

	var totalRequests = 0
	for _, v := range data {
		totalRequests += v
	}

	fmt.Printf("Total requests: %d\n", totalRequests)

	for k, v := range data {
		percentage := float64(v) / float64(totalRequests) * 100
		fmt.Printf("HTTP %s: %d requests (%0.2f %%)\n", k, v, percentage)
	}
}

func printIpReport(data map[string]int) {
	if len(data) == 0 {
		fmt.Println("the length of data is 0")
		return
	}

	fmt.Println("Unique addresses:", len(data))

	for k, v := range data {
		fmt.Printf("IP %s got %d hits.\n", k, v)
	}
}

func GenerateReport(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("\"%s\" no such file or directory", filepath)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	respCodeMap := make(map[string]int)
	ipAddrMap := make(map[string]int)

	startTime := time.Now()
	for scanner.Scan() {
		logEntry, err := parseLine(scanner.Text())

		if err != nil {
			fmt.Print(err)
			continue
		}

		ipAddrMap[logEntry.ip]++
		respCodeMap[logEntry.respCode]++
	}
	elapsed := time.Since(startTime)

	fmt.Println("Log source:", filepath)
	fmt.Println("Log parsing time:", elapsed)

	printRequestReport(respCodeMap)
	printIpReport(ipAddrMap)

	return nil
}
