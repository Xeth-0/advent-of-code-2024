package main

// Too many imports, ik.
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Vector struct {
	x int
	y int
}

type ClawMachine struct {
	ButtonA Vector
	ButtonB Vector
	Prize   Vector
}

func main() {
	// Parse args
	runExampleInputFlag := flag.Bool("e", true, "run the example input or the test input")
	loggingActiveFlag := flag.Bool("l", false, "logging active or not")

	flag.Parse()

	runExampleInput := *runExampleInputFlag
	loggingActive := *loggingActiveFlag

	// Set up logger
	if loggingActive || runExampleInput { // Set log file
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
	clawMachines := make([]ClawMachine, 0)
	for scanner.Scan() {
		// Going to read 3 lines at once, get all the data for the claw machine at once.
		line1 := strings.Split(strings.Split(scanner.Text(), ": ")[1], ", ")
		btnAParams := parseButton(line1)

		scanner.Scan()
		line2 := strings.Split(strings.Split(scanner.Text(), ": ")[1], ", ")
		btnBParams := parseButton(line2)

		scanner.Scan()
		line3 := strings.Split(strings.Split(scanner.Text(), ": ")[1], ", ")
		prize := parsePrize(line3)

		clawMachine := ClawMachine{
			ButtonA: btnAParams,
			ButtonB: btnBParams,
			Prize:   prize}

		clawMachines = append(clawMachines, clawMachine)

		scanner.Scan()
	}
	log.Println("Original Input: ", clawMachines)

	// Part One
	tokens := 0
	for _, clawMachine := range clawMachines {
		tokens += minimumTokensRequired(clawMachine, false)
	}
	fmt.Println("Part A: ", tokens)

	// Part two
	tokens = 0
	for _, clawMachine := range clawMachines {
		tokens += minimumTokensRequired(clawMachine, true)
	}
	fmt.Println("Part B: ", tokens)
}

func minimumTokensRequired(clawMachine ClawMachine, part2 bool) int {
	btnA := clawMachine.ButtonA
	btnB := clawMachine.ButtonB
	prize := clawMachine.Prize

	if part2 {
		prize.x += 10000000000000
		prize.y += 10000000000000
	}

	numB := ((btnA.x * prize.y) - (prize.x * btnA.y)) / ((btnB.y * btnA.x) - (btnB.x * btnA.y))
	numA := (prize.x - (numB * btnB.x)) / btnA.x

	// Check if the solution works. otherwise, no solution.
	if ((numA*clawMachine.ButtonA.x)+(numB*clawMachine.ButtonB.x) != prize.x) || (numA*clawMachine.ButtonA.y)+(numB*clawMachine.ButtonB.y) != prize.y {
		// No solution
		numA = 0
		numB = 0
	}

	return 3*numA + numB
}

func parseButton(str []string) Vector {
	_x := strings.Split(str[0], "+")[1]
	_y := strings.Split(str[1], "+")[1]

	x, err := strconv.Atoi(_x)
	if err != nil {
		log.Fatal("error parsing int: ", err)
	}

	y, err := strconv.Atoi(_y)
	if err != nil {
		log.Fatal("error parsing int: ", err)
	}

	return Vector{x: x, y: y}
}

func parsePrize(str []string) Vector {
	_x := strings.Split(str[0], "=")[1]
	_y := strings.Split(str[1], "=")[1]

	x, err := strconv.Atoi(_x)
	if err != nil {
		log.Fatal("error parsing int: ", err)
	}

	y, err := strconv.Atoi(_y)
	if err != nil {
		log.Fatal("error parsing int: ", err)
	}

	return Vector{x: x, y: y}
}
