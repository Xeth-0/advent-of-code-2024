package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// DISABLE LOGGER FOR THE ACTUAL TEST. Too large, will end up logging >50mbs of unintelligible maps.

	day6( "../input/input.example.txt", true)
	// day6( "../input/input.txt", false)
}

type Guard struct {
	position  Vector // (x,y) coordinates for the guard position
	direction string // | x => facing right | -x => facing left | y => facing up | -y => facing down |
}

type Vector struct {
	x int
	y int
}

func day6(filePath string, loggerActive bool) {
	// Set up logging
	logFile, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	
	// Reading the input
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("file not found")
	}
	scanner := bufio.NewScanner(file)

	givenMap := make([][]string, 0, 1024)
	guard := Guard{}
	for scanner.Scan() {
		line := scanner.Text()
		_guardPosition := strings.Index(line, "^")
		if _guardPosition != -1 {
			// Guard Position found
			guard.position.y = len(givenMap)
			guard.position.x = _guardPosition
			guard.direction = "y"
		}
		splitLine := strings.Split(line, "")
		givenMap = append(givenMap, splitLine)
	}

	// Part One
	if loggerActive {
		log.Println("------------------_PART ONE_--------------------------")
	}
	testingMap := deepCopy(givenMap)
	_, locationsVisited, _ := testGuardPath(testingMap, guard, loggerActive)
	fmt.Println("Part One: ", locationsVisited)

	// Part Two
	log.Println("------------------_PART TWO_--------------------------")
	loops := 0
	for i := range len(givenMap) {
		for j := range len(givenMap[0]) {
			// Attempt to place an obstacle
			isGuardPosition := givenMap[i][j] == "^" // guard starting position, can't place obstacle here.
			isObstacle := givenMap[i][j] == "#"      // no point to placing a new obstacle here.
			isOnPath := testingMap[i][j] == "."      // efficiency trick. making new obstacles on the current guard path narrows the search space.
			if isObstacle || isGuardPosition || isOnPath {
				continue
			}

			// Place an obstacle at location.
			newMap := deepCopy(givenMap)
			newMap[i][j] = "0" // different symbol for better visualization.

			// Test the new map for a loop
			loopExists, _, newMap := testGuardPath(newMap, guard, false)
			if loopExists {
				loops++
			}

			// Log the result
			if loggerActive && loopExists {
				log.Println("----------------------------------------------------------")
				log.Println("Loops?: ", loopExists)
				log.Println("Locations Visited: ", locationsVisited)
				logMap(newMap)
				log.Println()
			}
		}
	}

	fmt.Println("Part Two: ", loops)
}

func deepCopy(givenMap [][]string) (copiedMap [][]string) {
	copiedMap = make([][]string, 0, len(givenMap))

	for i, row := range givenMap {
		copiedMap = append(copiedMap, make([]string, 0, len(givenMap[i])))
		copiedMap[i] = append(copiedMap[i], row...)
	}

	return copiedMap
}

func testGuardPath(givenMap [][]string, guard Guard, loggerActive bool) (loopExists bool, locationsVisited int, updatedMap [][]string) {
	// Basically trying to make a set for tuples. since go doesn't have either, had to get stupid to make it work.
	visited := make(map[int](map[int][]string), 0)
	for i := range len(givenMap) {
		visited[i] = make(map[int][]string)
	}

	// let's mark the starting point as visited.
	// Visited will also store the direction that the guard was facing at that position as a list of directions: ["x", "y", "-x", ...] to detect loops.
	visited[guard.position.x][guard.position.y] = append(visited[guard.position.x][guard.position.y], guard.direction)
	locationsVisited = 1
	outOfBounds := false
	testingMap := givenMap

	// Let's map out the guard's path
	for {
		outOfBounds, testingMap = guard.move(testingMap)

		if loggerActive {
			log.Println("-----------------------------------------------------------")
			log.Println("New Guard Position: ", guard.position)
			logMap(givenMap)
		}

		// Check if the guard is still in the given bounds of the map.
		if outOfBounds {
			break
		}

		// Check previously visited positions.
		directions, exist := visited[guard.position.x][guard.position.y]
		if !exist {
			visited[guard.position.x][guard.position.y] = make([]string, 0, 1)
			directions = visited[guard.position.x][guard.position.y]
		}

		// Check for a loop. If the guard is at a position with a direction it's been before, it's in a loop.
		if len(directions) > 0 {
			for _, dir := range directions {
				if dir == guard.direction { // Loop found
					return true, locationsVisited, givenMap
				}
			}
		} else {
			// New Distinct Location. Increment the counter.
			locationsVisited++
		}

		// Mark the location and direction combo as visited.
		visited[guard.position.x][guard.position.y] = append(directions, guard.direction)

	}

	return false, locationsVisited, givenMap
}

func logMap(givenMap [][]string) {
	log.Println()
	for _, l := range givenMap {
		log.Println(l)
	}
	log.Println()
}

// Moves the guard, updates the map, returns outOfBounds? and the new map.
func (guard *Guard) move(givenMap [][]string) (outOfBounds bool, updatedMap [][]string) {
	// saving the old position
	oldPosition := guard.position
	oldDirection := guard.direction

	switch guard.direction {
	case "x":
		y := guard.position.y
		x := guard.position.x + 1
		if x < len(givenMap[0]) && (givenMap[y][x] == "#" || givenMap[y][x] == "0") {
			// execute a right turn
			guard.direction = "-y"
			givenMap[guard.position.y][guard.position.x] = "+"
			return guard.move(givenMap)
		}
		guard.position.x += 1

	case "-x":
		y := guard.position.y
		x := guard.position.x - 1
		if x >= 0 && (givenMap[y][x] == "#" || givenMap[y][x] == "0") {
			// execute a right turn
			guard.direction = "y"
			givenMap[guard.position.y][guard.position.x] = "+"
			return guard.move(givenMap)
		}
		guard.position.x -= 1

	case "y":
		y := guard.position.y - 1
		x := guard.position.x
		if y >= 0 && (givenMap[y][x] == "#" || givenMap[y][x] == "0") {
			// execute a right turn
			guard.direction = "x"
			givenMap[guard.position.y][guard.position.x] = "+"
			return guard.move(givenMap)
		}
		guard.position.y -= 1

	case "-y":
		y := guard.position.y + 1
		x := guard.position.x
		if y < len(givenMap) && (givenMap[y][x] == "#" || givenMap[y][x] == "0") {
			// execute a right turn
			guard.direction = "-x"
			givenMap[guard.position.y][guard.position.x] = "+"
			return guard.move(givenMap)
		}
		guard.position.y += 1
	}

	// Update the map. This is effectively useless, but let's me print the map to make a cool viz.
	// Let's update the old position to show the direction of mov't.
	oldDirectionIcon := givenMap[oldPosition.y][oldPosition.x] // Remove the guard indicator newPositionIcon from before.
	newDirectionIcon := oldDirectionIcon
	if oldDirectionIcon == "." {
		switch oldDirection {
		case "x", "-x":
			newDirectionIcon = "~"
		case "y", "-y":
			newDirectionIcon = "|"
		}
	} else if oldDirectionIcon == "|" || oldDirectionIcon == "-" {
		newDirectionIcon = "+"
	}
	givenMap[oldPosition.y][oldPosition.x] = newDirectionIcon

	// Check if the new guard position is within bounds of the map.
	outOfBoundsX := (guard.position.x < 0) || (guard.position.x >= len(givenMap[0]))
	outOfBoundsY := (guard.position.y < 0) || (guard.position.y >= len(givenMap))

	if outOfBoundsX || outOfBoundsY {
		// Now let's visualize the new position of the guard.
		newPositionIcon := "^"   // Icon for the new position
		switch guard.direction { // Make it reflect the direction the guard is facing. just for fun
		case "x":
			newPositionIcon = ">"
		case "y":
			newPositionIcon = "^"
		case "-x":
			newPositionIcon = "<"
		case "-y":
			newPositionIcon = "V"
		}
		givenMap[oldPosition.y][oldPosition.x] = newPositionIcon

		return true, givenMap
	}

	return false, givenMap
}
