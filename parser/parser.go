package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var logPattern = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[([^\]]+)\] "(\S+) ([^"]*) (\S+)" (\d{3}) (\S+) "([^"]*)" "([^"]*)" "([^"]*)"$`)

func parseLine(line string) (string, error) {
	m := logPattern.FindStringSubmatch(line)

	if m == nil {
		return "", fmt.Errorf("line does not match log pattern: %q\n", string(line))
	}

	// index 8 is HTTP status code in regex
	respCode := m[8]
	return respCode, nil
}

func printReport(data map[string]int) {
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

func GenerateReport(filepath string) error {
	respCodeMap := make(map[string]int)

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("\"%s\" no such file or directory", filepath)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		respCode, err := parseLine(scanner.Text())

		if err != nil {
			log.Info(err)
			continue
		}

		respCodeMap[respCode]++
	}

	printReport(respCodeMap)
	return nil
}
