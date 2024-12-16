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

type Vector struct {
	x int
	y int
}

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

	// Read the map.
	givenMap := make([][]string, 0, 1024)
	robotPos := Vector{}
	robotFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// end of map section
			break
		}
		splitLine := strings.Split(line, "")

		// Find the robot position, if not already tracked
		if !robotFound {
			for i, char := range splitLine {
				if char == "@" {
					robotPos.y = len(givenMap)
					robotPos.x = i
				}
			}
		}
		
		givenMap = append(givenMap, splitLine)

	}
	log.Println("Original Input: ")
	plotMap(givenMap)

	// Read the robot's movements
	rawRobotMoves := ""
	for scanner.Scan() {
		rawRobotMoves += scanner.Text()
	}
	robotMoves := strings.Split(rawRobotMoves, "")

	log.Println("Robot Moves:")
	log.Println(rawRobotMoves)

	// Part One
	moveRobot(givenMap, robotMoves)
}

func moveRobot(givenMap [][]string, robotMoves []string, robotPos Vector) {
	// Needs to be able to shove boxes when there is no wall behind them
	// Let's track the robot's position
	for _, m := range robotMoves {
		move := robotPos.add(getMove(m))
		fmt.Println(move)
		break
		// moveObstacle(&givenMap, ,move)
	}
}

func moveObstacle(givenMap *[][]string, coord Vector, move Vector) {

}

func (vec Vector) add(vec2 Vector) Vector {
	vec.x += vec2.x
	vec.y += vec2.y
	return vec
}

func getMove(move string) Vector {
	robotMove := Vector{}
	switch move {
	case "^":
		robotMove.y = 1
	case "v":
		robotMove.y = -1
	case ">":
		robotMove.x = 1
	case "<":
		robotMove.x = -1
	}

	return robotMove
}

func plotMap(givenMap [][]string) {
	for _, row := range givenMap {
		log.Println(row)
	}
}
