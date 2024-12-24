package main

// Too many imports, ik.
import (
	"advent2024/utils/vector"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Define controls.
var NUMPAD = make(map[string]vector.Vector)
var DPAD = make(map[string]vector.Vector)

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
	// Read the input
	codes := readInput(filePath)

	// Setup numpad and dpad
	var nums = "_0A123456789"
	var directions = "<v>_^A"

	for i, char := range strings.Split(nums, "") {
		NUMPAD[char] = vector.Vector{X: i % 3, Y: i / 3}
	}
	for i, char := range strings.Split(directions, "") {
		DPAD[char] = vector.Vector{X: i % 3, Y: i / 3}
	}

	// Run the solver.
	for _, part := range parts {
		switch part {
		case 1:
			clearMemo()
			solve(codes, 2)
		case 2:
			clearMemo()
			solve2(codes, 25)
		}
	}
}

// Task: control a robot on a keypad, that's controlled by a robot, that's controlled by a robot, that's controlled by me.
//   Find input => Robot 1 => Robot 2 => Robot 3 => Desired Numbers.
//   so basically i need 4 transformations to an output. and given that output, i need to find the input to make it happen.
// Approach: Work backwards. Determine the moves required to hit the final keypad, then determine the
//   inputs from the directional dpad to get those moves, and repeat.

// MOVEMENT PRECEDENCE: RIGHT - UP - DOWN - LEFT => to prevent robot arm from accidentally hovering over the empty key.

// Map the numpad and dpad to Vectors(x,y). This'll make getting the moves required easier.

// Depth => number of robots with the dpad (ignoring the final one with the numpad).
func solve(codes []string, depth int) {

	// fmt.Println("Numpad: ", numpad)
	// fmt.Println("Dpad: ", dpad)

	requiredInput := make(map[string][]string)
	for _, code := range codes {
		// Determine the input required at each stage of the robots(working backwards)
		fmt.Println(code)
		moves := getInput(code, depth+1, true)

		requiredInput[code] = append(requiredInput[code], moves)
	}

	finalComplexity := 0
	fmt.Println("Part One: ")
	// fmt.Println("Required Inputs: ", requiredInput)
	for code, moves := range requiredInput {
		finalMove := moves[len(moves)-1]
		complexity := complexityScore(code, len(finalMove))
		finalComplexity += complexity
		fmt.Printf("%s: %d - %s\n", code, complexity, finalMove)
	}

	fmt.Println("Final Complexity: ", finalComplexity)
}

func solve2(codes []string, depth int) {
	// Solver for part two.
	requiredInput := make(map[string]int)
	for _, code := range codes {
		// Determine the input required at each stage of the robots(working backwards)
		fmt.Println(code)
		moves := getInput2(code, depth+1, true)

		requiredInput[code] = moves
	}

	finalComplexity := 0
	fmt.Println("Part One: ")
	// fmt.Println("Required Inputs: ", requiredInput)
	for code, moves := range requiredInput {
		finalMove := moves
		complexity := complexityScore(code, finalMove)
		finalComplexity += complexity
		fmt.Printf("%s: %d - %d\n", code, complexity, finalMove)
	}

	fmt.Println("Final Complexity: ", finalComplexity)
}


// For a given code, returns the commands that need to be sent from the pupeteer to make this robot press that code.
func getInput(code string, depth int, start bool) string {
	// I need to get the most optimal sequence of moves, at the lowest depth (the input level),
	// that produces the given code at the highest level.
	// So this is kind of like a search problem, find the most optimal path. Except, I myself can't
	// figure out the most optimal path. and the optimal path for one level might not be optimal a layer
	// deeper. So instead, i'm going to try to generate all possible inputs at depth 0 for the given code
	// at depth n, and choose the shortest. so, recursion!!!!, but i should memoize this, or it'll get ridiculous.
	// approach: 2 possible moves at each depth: horizontal first || vertical first.
	// going to send both down the recursive hell, and choose the best moves
	// i want to minimize at the point that all layers below press A. so that'll be each character.
	
	if depth <= 0 {
		return code
	}

	var controls map[string]vector.Vector
	if start {
		controls = NUMPAD
	} else {
		controls = DPAD
	}

	move := ""
	cursor := controls["A"] // Start the cursor at "A"

	for _, char := range code {
		target := controls[string(char)]
		numMoves, moves := getRoute(cursor, target, controls["_"]) // Translate the code into the two possible options

		if numMoves == 2 {
			newMoves := []string{getInput(moves[0], depth-1, false), getInput(moves[1], depth-1, false)} // Recursively find the most optimum move to achieve this at max depth.
			if len(newMoves[0]) < len(newMoves[1]) {
				move += newMoves[0]
			} else {
				move += newMoves[1]
			}
		} else {
			// Just use the best move
			move += getInput(moves[0], depth-1, false)
		}

		// Move the cursor
		cursor = controls[string(char)]
	}
	return move
}


var memo map[string]int = make(map[string]int)

func getInput2(code string, depth int, start bool) int {
	// For part two
	// Now this is the same as the first getInput function. But since we're doing depth 25, plain recursion is not gonna cut it.
	// Need to memoize. And since the moves are too long, and depend on depth and the code, it keeps blowing out the virtual memory.
	// Just going to return length because i'm kinda tired and there's other shit to do.
	// It's bonkers how much faster this is, first shot solution is correct too.
	if depth <= 0 {
		return len(code)
	}

	memo_key := fmt.Sprintf("%s_%d", code, depth)
	if val, cached := memo[memo_key]; cached {
		return val
	}

	var controls map[string]vector.Vector
	if start {
		controls = NUMPAD
	} else {
		controls = DPAD
	}

	move := 0
	cursor := controls["A"] // Start the cursor at "A"

	for _, char := range code {
		target := controls[string(char)]
		numMoves, moves := getRoute(cursor, target, controls["_"]) // Translate the code into the two possible options

		if numMoves == 2 {
			newMoves := []int{getInput2(moves[0], depth-1, false), getInput2(moves[1], depth-1, false)} // Recursively find the most optimum move to achieve this at max depth.
			if newMoves[0] < newMoves[1] {
				move += newMoves[0]
			} else {
				move += newMoves[1]
			}
		} else {
			// Just use the best move
			move += getInput2(moves[0], depth-1, false)
		}

		// Move the cursor
		cursor = controls[string(char)]
	}
	memo[memo_key] = move
	return move
}

func clearMemo() {
	memo = make(map[string]int)
}


// This function just returns the two options, with one constraint. None of them pass over the bad tile.
func getRoute(start vector.Vector, end vector.Vector, badTile vector.Vector) (numRoutes int, routes [2]string) {
	// All this function cares about, is a start position and end position. It'll return the
	// two (yes, two), options of moves that result in the cursor moving from start => end.
	// Now, there are two ways of constructing a solution, that are more efficient than the others.
	// And that's to move horizontally, then vertically. Or to move vertically, or horizontally.
	// Any sequence of moves that interleaves these two, will not be efficient as it moves up a layer to the
	// robot controlling this robot. let's say the moves are >v>v. Robot controlling this one will have to
	// make this robot navigate to >, then v, then >, and then back to v. Avoiding interleaving would've made this
	// just two navigations to > then v or vice versa.

	// Distance we need to move in each axes.
	displacement := vector.Vector{
		X: end.X - start.X,
		Y: end.Y - start.Y,
	}

	// Construct the horizontal and vertical moves
	horizontal := strings.Repeat(">", max(0, displacement.X)) + strings.Repeat("<", max(-displacement.X, 0))
	vertical := strings.Repeat("^", max(0, displacement.Y)) + strings.Repeat("v", max(-displacement.Y, 0))

	// Return value
	options := [2]string{}

	if horizontal == "" {
		options[0] = vertical + "A"
		return 1, options
	}
	if vertical == "" {
		options[0] = horizontal + "A"
		return 1, options
	}

	// Now, to test for the bad tile.

	if badTile.X == end.X && badTile.Y == start.Y { // so after moving horizontally, but before vertically, we're at badtile
		// Cannot move horizontally first
		options[0] = vertical + horizontal + "A"
		options[1] = ""
		return 1, options
	} else if end.Y == badTile.Y && start.X == badTile.X { // after moving vertically, but before horizontally, we're at badtile
		// Cannot move vertically first
		options[0] = horizontal + vertical + "A"
		options[1] = ""
		return 1, options
	}

	options[0] = horizontal + vertical + "A"
	options[1] = vertical + horizontal + "A"

	return 2, options
}

func complexityScore(code string, sequence int) int {
	val, err := strconv.Atoi(code[:3])
	if err != nil {
		log.Fatal("error parsing the numeric part of code: ", err)
	}

	return sequence * val
}

func readInput(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([]string, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
	}

	return buffer
}
