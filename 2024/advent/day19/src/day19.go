package main

// Too many imports, ik.
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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

	day12(filePath, 1, 2)
}

func day12(filePath string, parts ...int) {
	designs, patterns := readInput(filePath)
	for _, part := range parts {
		switch part {
		case 1:
			part1(designs, patterns)
		case 2:
			part2(designs, patterns)
		}
	}
}

func part1(designs []string, patterns []string) {
	possibleDesigns := 0
	for _, design := range designs {
		if (countWays(design, &patterns, &map[string]bool{})) {
			possibleDesigns++
		}
	}

	fmt.Println("Part One: ", possibleDesigns)
}

func countWays(targetDesign string, patterns *[]string, memo *map[string]bool) bool{
	if targetDesign == ""{
		return true
	}

	// check the memo 
	if ((*memo)[targetDesign]) {
		return (*memo)[targetDesign]
	}

	for _, pattern := range (*patterns) {
		// check if the pattern matches with the start of the target design
		if (len(pattern) <= len(targetDesign) && pattern == targetDesign[:len(pattern)]) {
			if (countWays(targetDesign[len(pattern):], patterns, memo)) {
				return true
			}
		}
	}
	return false
}

func part2(designs []string, patterns []string) {
	ways := 0 
	for _, design := range designs {
		ways += countWays2(design, &patterns, &map[string]int{})
	}

	fmt.Println("Part Two: ", ways)
}

func countWays2(targetDesign string, patterns *[]string, memo *map[string]int) int {
	if targetDesign == ""{
		return 1
	}

	// check the memo
	if ((*memo)[targetDesign] > 0) {
		return (*memo)[targetDesign]
	}

	// count the number of ways the design can be constructed
	totalWays := 0
	for _, pattern := range (*patterns) {
		// check if the pattern matches with the start of the target design
		if (len(pattern) <= len(targetDesign) && pattern == targetDesign[:len(pattern)]) {
			totalWays += countWays2(targetDesign[len(pattern):], patterns, memo)
		}
	}
	(*memo)[targetDesign] = totalWays
	return totalWays
}

func readInput(filePath string) (designs []string, patterns []string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	// Read the desings (single line)
	scanner.Scan()
	line := strings.Split(scanner.Text(), ", ")
	fmt.Println("Designs: ", line)
	patterns = line

	// Skip the empty line
	scanner.Scan()

	// Read the patterns to match
	for scanner.Scan() {
		pattern := scanner.Text()
		designs = append(designs, pattern)
	}
	fmt.Println("Patterns: ", patterns)

	return designs, patterns
}
