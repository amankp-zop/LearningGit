package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// This functions aims to read logs and perform analysis.
func main() {
	// Check if a filename is passed
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go fileName")
		return
	}

	filename := os.Args[1]

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close() // ensure the file is closed when the function ends

	// Counters
	var infoCount, warningCount, errorCount, totalCount int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		totalCount++

		switch {
		case strings.HasPrefix(line, "[INFO]"):
			infoCount++
		case strings.HasPrefix(line, "[WARNING]"):
			warningCount++
		case strings.HasPrefix(line, "[ERROR]"):
			errorCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Display summary
	fmt.Printf("Log Analysis of file: %s\n\n", filename)
	fmt.Printf("INFO: %d entries (%.2f%%)\n", infoCount, percent(infoCount, totalCount))
	fmt.Printf("WARNING: %d entries (%.2f%%)\n", warningCount, percent(warningCount, totalCount))
	fmt.Printf("ERROR: %d entries (%.2f%%)\n", errorCount, percent(errorCount, totalCount))

	// Timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("\nAnalyzed at: %s\n", currentTime)
}

// percent returns percentage of count over total
func percent(count, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(count) / float64(total) * 100
}
