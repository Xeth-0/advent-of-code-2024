package main

// Too many imports, ik.
import (
	"os"
	"fmt"
	"log"
	"flag"
	"bufio"
	"strconv"
	"strings"
)

func main() {
	// Parse args
	runExampleInputFlag := flag.Bool("e", true, "run the example input or the test input")
	loggingActiveFlag := flag.Bool("l", false, "logging active or not")

	flag.Parse()

	runExampleInput := *runExampleInputFlag
	loggingActive := *loggingActiveFlag

	// Set up logger
	if loggingActive || runExampleInput {
		// Set log file
		logFile, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}

	var filePath string
	if !runExampleInput {
		fmt.Println("Running Test input...")
		filePath = "../input/input.txt"
	} else {
		fmt.Println("Running Example input...")
		filePath = "../input/input.example.txt"
	}

	day12(filePath)
}

func day12(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([]int, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		for _, v := range line {
			val, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("error parsing int: %v", err)
			}

			buffer = append(buffer, val)
		}
	}
	log.Println("Original Input: ", buffer)

}
