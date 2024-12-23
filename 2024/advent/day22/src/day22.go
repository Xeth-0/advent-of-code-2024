package main

// Too many imports, ik.
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
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
	secretNumbers := readInput(filePath)
	solve(secretNumbers)
}

func solve(secretNumbers []int) {
	// original := make([]int, 0)
	// original = append(original, secretNumbers...)

	windows := make(map[[4]int]int)

	for i, secretNumber := range secretNumbers {
		last5 := [5]int{}
		last5[0] = secretNumber % 10

		// To keep track of the first time each window of 4 occurs. Monkey sells at first chance, so latter ones should be ignored.
		visited := make(map[[4]int]bool)

		for f := range 2000 {
			nextSecretNumber := prune(mix(secretNumbers[i]*64, secretNumbers[i]))
			nextSecretNumber = prune(mix(nextSecretNumber/32, nextSecretNumber))
			nextSecretNumber = prune(mix(nextSecretNumber*2048, nextSecretNumber))
			secretNumbers[i] = nextSecretNumber

			if f <= 3 { // for the first 5 secret numbers (including the first)
				last5[f+1] = nextSecretNumber % 10
				// fmt.Println("last5: ", last5)
			} else {
				// shift/re-adjust the window
				last5 = shiftLeft(last5)
				last5[4] = nextSecretNumber % 10
			}

			// Find the 4 consecutive changes, only if there have been 4 new secret numbers.
			if f >= 3 {
				window := findChange(last5)

				if !visited[window] {
					windows[window] += last5[4]
					visited[window] = true
				}
			}
		}
	}

	// Part One
	s := 0
	for _, secretNumber := range secretNumbers {
		// fmt.Println(original[i], ": ", secretNumber)
		s += secretNumber
	}
	fmt.Println("Part One - Sum: ", s)

	// Part Two
	maxVal := math.Inf(-1)
	bestWindow := [4]int{}
	for window, val := range windows {
		if val > int(maxVal) {
			bestWindow = window
			maxVal = float64(val)
		}
	}

	fmt.Println("Part Two - Best Window: ", bestWindow, "Bananas: ", maxVal)
}

func findChange(last5 [5]int) [4]int {
	change := [4]int{}
	for i := 1; i < 5; i++ {
		change[i-1] = last5[i] - last5[i-1]
	}

	return change
}

func shiftLeft(arr [5]int) [5]int {
	for i := 1; i < 5; i++ {
		arr[i-1] = arr[i]
	}
	arr[4] = 0
	return arr
}

func mix(secretNumber int, val int) int {
	return val ^ secretNumber
}

func prune(secretNumber int) int {
	return secretNumber % 16777216
}

func readInput(filePath string) []int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([]int, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := scanner.Text()
		val, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("error parsing int: %v", err)
		}

		buffer = append(buffer, val)

	}
	return buffer
}
