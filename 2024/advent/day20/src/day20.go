package main

// Too many imports, ik.
import (
	tui "advent2024/utils/tui"
	. "advent2024/utils/vector"
	"math"
	"time"

	// "time"

	tea "github.com/charmbracelet/bubbletea"

	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	WALL  = "#"
	START = "S"
	END   = "E"
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

	day12(filePath, runExampleInput, 1, 2)
}

func day12(filePath string, isExample bool, parts ...int) {
	grid := readInput(filePath)
	for _, part := range parts {
		switch part {
		case 1:
			// part1(grid, 100, isExample)
		case 2:
			part2(grid, 100, 20)
		}
	}
}

func part1(grid [][]string, distanceToSave int, isExample bool) {
	// Task: Count the number of paths, where a cheat can be used that can save 100 picoseconds(steps)
	// Approach: There is one actual path. The cheat will allow to skip thru 2 positions(even if theyre walls),
	//   and it will count if it saves >= 100 steps. This also means the start position of the cheat and the end
	//   position are on the original path. I can trace thru the original racetrack, check if other positions are
	//   separated from that position by a wall, and check the new cost?
	p := tea.NewProgram(tui.NewModel("Day 20 - Race Conditions"), tea.WithAltScreen(), tea.WithMouseCellMotion())

	go func() {
		path := findPath(grid)
		for _, vec := range path {
			grid[vec.Y][vec.X] = "0"
		}

		// let's get their distances (and since the path is sequential, ...)
		distances := make(map[Vector]int)
		for i, v := range path {
			distances[v] = len(path) - i - 1
		}

		deltas := []Vector{ // To skip once in the 4 directions. Could make it two, this is just for the viz.
			{X: 0, Y: 1},
			{X: 0, Y: -1},
			{X: 1, Y: 0},
			{X: -1, Y: 0},
		}

		// Count Cheats that skip at least the given number of tiles.
		viableCheats := 0
		for _, node := range path {
			// for each node, check it's 4 neighbors.
			for _, delta := range deltas {
				skip1 := node.Add(delta) // First tile skipped.
				dest := skip1.Add(delta) // Where the racecar will land after the cheat.

				if _, exists := distances[dest]; !exists || !dest.InBounds(len(grid[0]), len(grid)) {
					continue
				}

				// Update the viz. All this code is pretty unecessary but...
				temp1 := grid[skip1.Y][skip1.X]
				temp2 := grid[dest.Y][dest.X]

				grid[skip1.Y][skip1.X] = tui.BoxStyle.Render("^")
				grid[dest.Y][dest.X] = tui.BoxStyle.Render("@")

				if isExample {
					p.Send(tui.UpdateViewport("\n"+plotMap(grid), 100))
				}

				// The distance the skip has saved.
				distanceDelta := distances[node] - distances[dest] - 2
				if distanceDelta >= distanceToSave {
					viableCheats++

					if isExample { // the test takes too long otherwise.
						p.Send(tui.UpdateViewport(fmt.Sprintf("Skip Found: Saves %d steps\n%s", distanceDelta, plotMap(grid)), 100))
						time.Sleep(10 * time.Millisecond)
					}
				}

				// Reset the grid.
				grid[skip1.Y][skip1.X] = temp1
				grid[dest.Y][dest.X] = temp2
			}
		}

		// Display the result
		p.Send(tui.UpdateViewport(fmt.Sprintf("Cheats: %d\n%s", viableCheats, plotMap(grid)), 100))
		log.Println("Part One: ", viableCheats)
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal("error running tui: ", err)
	}
}

func part2(grid [][]string, distanceToSave int, cheatSpan int) {
	// Slightly different approach for this one. Since now it's upto 20 steps, and not limited to exactly 20 steps:
	//    => for each node, subtract 100 from it's distance and find any places on track that are within that distance to the end and
	//      within 20 steps from the node. pretty straight forward.
	//    => puzzle also has the constraint that cheats have to be unique(unique start and end), but this takes care of that too.

	// going to skip the viz for this as well.

	// Find the original path
	path := findPath(grid)
	for _, vec := range path {
		grid[vec.Y][vec.X] = "0"
	}

	// let's get their distances to the end.
	distances := make(map[Vector]int)
	for i, v := range path {
		distances[v] = len(path) - i - 1
	}

	// Now, this will be O(n^2). I wish it was better. But welp
	viableCheats := 0
	for i, startNode := range path {
		for _, destNode := range path[i:] {
			deltaX := float64(startNode.X - destNode.X)
			deltaY := float64(startNode.Y - destNode.Y)
			cheat := int(math.Abs(deltaX) + math.Abs(deltaY))

			distanceSaved := distances[startNode] - distances[destNode] - cheat

			if cheat <= cheatSpan && distanceSaved >= distanceToSave {
				viableCheats++
			}
		}
	}

	fmt.Println("Part Two: ", viableCheats)
}

// Assumes there is one valid path, and it's the only path. Just because that's how the puzzle input works.
func findPath(grid [][]string) []Vector {
	startingPosition := Vector{}

loop:
	for r, row := range grid {
		for c := range row {
			if grid[r][c] == "S" {
				startingPosition.X = c
				startingPosition.Y = r

				break loop
			}
		}
	}

	// Trace the path. i'm kinda dumb, so just gonna do the most obvious bfs around
	path := make([]Vector, 0, 1)
	queue := []Vector{startingPosition}
	visited := make(map[Vector]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// add it to the path. because it's obviously on it.
		visited[current] = true
		path = append(path, current)

		if grid[current.Y][current.X] == "E" {
			break
		}

		for _, neighbor := range current.Neighbors() {
			tile := grid[neighbor.Y][neighbor.X]
			if neighbor.InBounds(len(grid[0]), len(grid)) && !visited[neighbor] && (tile == "." || tile == "E") {
				queue = append(queue, neighbor)
				break
			}
		}
	}

	return path
}

func readInput(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	grid := make([][]string, 0, 1024)

	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "")
		grid = append(grid, row)
	}

	// log.Println("Original Input: ", grid)
	return grid
}

func plotMap(grid [][]string) string {
	// strings.Builder is more efficient than simple concatenation, massive speedup using this instead.
	var builder strings.Builder

	for _, row := range grid {
		for _, c := range row {
			switch c {
			case WALL:
				builder.WriteString(tui.WallStyle.Render(WALL))
			default:
				builder.WriteString(c)
			}
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}
