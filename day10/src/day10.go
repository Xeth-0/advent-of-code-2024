package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Parse args
	runExampleInputFlag := flag.Bool("e", true, "run the example input or the test input")
	flag.Parse()

	runExampleInput := *runExampleInputFlag

	// Set up logger
	if runExampleInput {
		// Set log file
		logFile, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}
	log.SetFlags(2) // disable timestamp

	var filePath string
	if !runExampleInput {
		fmt.Println("Running Test input...")
		filePath = "../input/input.txt"
	} else {
		fmt.Println("Running Example input...")
		filePath = "../input/input.example.txt"
	}

	day10(filePath)
}

type Vector struct {
	xPos int
	yPos int
}

func day10(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([][]int, 0, 1024)

	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		parsedLine := make([]int, 0, len(line))
		for _, v := range line {
			val, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("error parsing int: %v", err)
			}

			parsedLine = append(parsedLine, val)
		}

		buffer = append(buffer, parsedLine)
	}

	log.Println("Original Input:")
	logMap(buffer)

	// Find 0's
	trailHeads := make([][]Vector, 0)
	for i, row := range buffer {
		for j, val := range row {
			if val == 0 {
				t := trailBlazing(buffer, Vector{
					xPos: j,
					yPos: i,
				})
				if len(t) > 0 {
					trailHeads = append(trailHeads, t)
				}
			}
		}
	}

	count := countUniqueTrailHeads(trailHeads)
	fmt.Println("PART ONE: ", count)

	count = computeRating(trailHeads)
	fmt.Println("PART TWO: ", count)
}

// Returns the whatever it needs
func trailBlazing(givenMap [][]int, startingPos Vector) []Vector {
	if givenMap[startingPos.yPos][startingPos.xPos] == 9 {
		return []Vector{startingPos}
	}

	trailHeads := make([]Vector, 0) // Stores the locations of trail-heads (9)
	currentHeight := givenMap[startingPos.yPos][startingPos.xPos]

	if startingPos.yPos+1 < len(givenMap) && givenMap[startingPos.yPos+1][startingPos.xPos] == currentHeight+1 {
		t := trailBlazing(givenMap, Vector{
			xPos: startingPos.xPos,
			yPos: startingPos.yPos + 1,
		})
		if len(t) > 0{
			trailHeads = append(trailHeads, t...)
		}
	}
	if startingPos.yPos-1 >= 0 && givenMap[startingPos.yPos-1][startingPos.xPos] == currentHeight+1 {
		t := trailBlazing(givenMap, Vector{
			xPos: startingPos.xPos,
			yPos: startingPos.yPos - 1,
		})
		if len(t) > 0{
			trailHeads = append(trailHeads, t...)
		}
	}
	if startingPos.xPos+1 < len(givenMap[0]) && givenMap[startingPos.yPos][startingPos.xPos+1] == currentHeight+1 {
		t := trailBlazing(givenMap, Vector{
			xPos: startingPos.xPos + 1,
			yPos: startingPos.yPos,
		})
		if len(t) > 0{
			trailHeads = append(trailHeads, t...)
		}
	}
	if startingPos.xPos-1 >= 0 && givenMap[startingPos.yPos][startingPos.xPos-1] == currentHeight+1 {
		t := trailBlazing(givenMap, Vector{
			xPos: startingPos.xPos - 1,
			yPos: startingPos.yPos,
		})
		if len(t) > 0{
			trailHeads = append(trailHeads, t...)
		}
	}

	return trailHeads
}

func logMap(givenMap [][]int) {
	log.Println()
	for _, l := range givenMap {
		log.Println(l)
	}
	log.Println()
}

func countUniqueTrailHeads(trailHeads [][]Vector) int {
	count := 0
	for _, trail := range trailHeads {
		// for each trail (starting from a 0), count the number of unique trail heads
		uniqueTrailHeads := make(map[Vector]bool)

		for _, trailHead := range trail {
			uniqueTrailHeads[trailHead] = true
		}

		count += len(uniqueTrailHeads)
	}

	return count
}

func computeRating(trailHeads [][]Vector) int {
	count := 0
	for _, trail := range trailHeads {
		count += len(trail)
	}
	return count
}