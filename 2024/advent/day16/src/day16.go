package main

// Too many imports, ik.
import (
	"advent2024/utils/tui"
	. "advent2024/utils/vector"
	"math"
	"time"

	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	empty      = "."
	wall       = "#"
	robotUp    = "^"
	robotDown  = "v"
	robotLeft  = "<"
	robotRight = ">"
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
		// log.SetOutput(logFile)
	}
	log.SetFlags(0)

	var filePath string
	if !runExampleInput {
		fmt.Println("Running Test input...")
		filePath = "../input/input.txt"
	} else {
		fmt.Println("Running Example input...")
		filePath = "../input/input.example.txt"
	}

	day16(filePath, runExampleInput)
}

func day16(filePath string, isExampleInput bool) {
	grid, startPos, endPos := readInput(filePath)
	solve(grid, startPos, endPos, isExampleInput)
}

func solve(grid [][]string, startPos Vector, endPos Vector, isExampleInput bool) {
	p := tea.NewProgram( // The terminal UI to draw the path and visualize the search.
		tui.Model{}, tea.WithAltScreen(), tea.WithMouseCellMotion(),
	)

	go func() {
		// Going to use UCS. Thought about A*, but it seems like UCS would be the better choice here.
		paths := make([]*Node, 0)
		bestPathCost := math.MaxInt

		startNode := Node{pos: startPos, cost: 0, dir: robotRight} // starts facing east.

		// Initialize the priority queue
		pq := make(PriorityQueue, 0)
		heap.Init(&pq)
		heap.Push(&pq, &startNode)

		// Set to track visited positions
		visited := make(map[Vector]map[string]bool)

		for pq.Len() > 0 {
			// this is just a type assertion. we're just casting the return value of Pop(an interface{}/generic), to a *Node pointer
			current := heap.Pop(&pq).(*Node)

			grid[current.pos.Y][current.pos.X] = current.dir // For the visualization, adding a cursor to the robot to keep track.

			// Update the visualizer. This will heavily slow down the program, so it's only going to run for the example input.
			if isExampleInput {
				s := fmt.Sprintf("Grid\nCost: %d\n", current.cost) + plotMap(grid)
				p.Send(tui.UpdateViewport(s, 1000))
				time.Sleep(10 * time.Millisecond) // just so i can keep up and see what's happening
			}

			// If the current path cost is greater than the best one so far, any path we find after this will not be a "best path". break.
			if current.cost > bestPathCost {
				break
			}

			// Check if we reached the goal.
			if current.pos.Equals(endPos) {
				bestPathCost = current.cost    // Also going to need the other paths with the same cost. So saving this.
				paths = append(paths, current) // Add the end node to the paths array. Can use node.parent to retrace steps.
			}

			// Mark the current node as visited. Takes direction into account.
			if _, exists := visited[current.pos]; !exists {
				visited[current.pos] = make(map[string]bool)
			}
			visited[current.pos][current.dir] = true

			// Grab the neighbors
			for _, neighbor := range current.Neighbors() {
				if neighbor.pos.InBounds(len(grid[0]), len(grid)) && grid[neighbor.pos.Y][neighbor.pos.X] != wall && !visited[neighbor.pos][neighbor.dir] {
					// neighbor.cost += endPos.Distance(neighbor.pos) // A* search. Will not always give the shortest path.
					neighbor.parent = current
					heap.Push(&pq, neighbor)
				}
			}

			grid[current.pos.Y][current.pos.X] = empty // For the visualization
		}

		// Draw the path, and count the number of nodes that are at least on a "best path"(soln to part two)
		nodesOnBestPath := 0
		onBestPath := make(map[Vector]bool)
		for _, node := range paths {
			for node.parent != nil {
				if !onBestPath[node.pos] {
					nodesOnBestPath++
					onBestPath[node.pos] = true
				}
				grid[node.pos.Y][node.pos.X] = tui.BoxStyle.Render("0")
				node = node.parent

				if isExampleInput { // Again, very slow on the test input. better to just render the final version.
					p.Send(tui.UpdateViewport("Tagging nodes on best paths...\n"+plotMap(grid), 100))
					time.Sleep(100 * time.Millisecond)
				}
			}
		}

		// Formatting the tui output.
		s := fmt.Sprintf("Best Path Cost: %d\n", bestPathCost)
		s += fmt.Sprintf("Nodes on best paths: %d\n", nodesOnBestPath)
		s += plotMap(grid)
		s+="\nPress 'q' or 'ctrl+c' or 'esc' to exit..."

		p.Send(tui.UpdateViewport(s, 2*len(grid[0]))) // Update the terminal ui.
	}()

	// Run the terminal UI. has to be in a separate thread.
	if _, err := p.Run(); err != nil {
		fmt.Println("error running visualizer: ", err)
		os.Exit(1)
	}
}

func readInput(filePath string) ([][]string, Vector, Vector) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([][]string, 0, 1024)
	startPos := Vector{}
	endPos := Vector{}

	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		for j, v := range line {
			if v == "S" {
				startPos.Set(j, len(buffer))
			} else if v == "E" {
				endPos.Set(j, len(buffer))
			}
		}
		buffer = append(buffer, line)
	}
	return buffer, startPos, endPos
}

func plotMap(givenMap [][]string) string {
	s := ""
	for _, row := range givenMap {
		for _, c := range row {
			switch c {
			case wall:
				s += tui.WallStyle.Render(wall)
			case robotUp:
				s += tui.RobotStyle.Render(robotUp)
			case robotLeft:
				s += tui.RobotStyle.Render(robotLeft)
			case robotDown:
				s += tui.RobotStyle.Render(robotDown)
			case robotRight:
				s += tui.RobotStyle.Render(robotRight)
			default:
				s += c
			}
		}
		s += "\n"
	}
	return s
}

// Returns the immediate the robot can visit. Direction matters. (Can only rotate once counter/clockwise as well)
func (node *Node) Neighbors() []*Node {
	// three neighbors at any position. Keep going same dir, rotate clockwise, rotate counterclockwise.

	// Could this have been cleaner, maybe. But i'm kinda dumb so...
	neighbors := make([]*Node, 0, 3)
	switch node.dir {
	case robotUp:
		neighbors = append(neighbors,
			&Node{
				pos:  node.pos.Up(),
				cost: node.cost + 1,
				dir:  node.dir,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotRight,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotLeft,
			},
		)
	case robotRight:
		neighbors = append(neighbors,
			&Node{
				pos:  node.pos.Right(),
				cost: node.cost + 1,
				dir:  node.dir,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotUp,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotDown,
			},
		)

	case robotDown:
		neighbors = append(neighbors,
			&Node{
				pos:  node.pos.Down(),
				cost: node.cost + 1,
				dir:  node.dir,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotRight,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotLeft,
			},
		)

	case robotLeft:
		neighbors = append(neighbors,
			&Node{
				pos:  node.pos.Left(),
				cost: node.cost + 1,
				dir:  node.dir,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotUp,
			},
			&Node{
				pos:  node.pos,
				cost: node.cost + 1000,
				dir:  robotDown,
			},
		)
	}
	return neighbors
}
