package main

// Too many imports, ik.
import (
	tui "advent2024/utils/tui"
	. "advent2024/utils/vector" // Yes, ik. i'm sorry. i won't change this.
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	bounds       Vector     // Bounds of the grid (also serve as the destination for the puzzle input)
	fallingBytes []Vector   // The bytes that are falling from the puzzle input (in order)
	board        [][]string // The board that is being plotted
}

type Node struct {
	pos    Vector
	parent *Node
}

func main() {
	// Parse args
	runExampleInputFlag := flag.Bool("e", true, "run the example input or the test input")
	loggingActiveFlag := flag.Bool("l", false, "logging active or not")

	flag.Parse()

	runExampleInput := *runExampleInputFlag
	_ = *loggingActiveFlag

	// Set up logger
	// Might need it even for the test case, since i can't log to console due to the tui.
	logFile, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	var filePath string
	if !runExampleInput {
		fmt.Println("Running Test input...")
		filePath = "../input/input.txt"
	} else {
		fmt.Println("Running Example input...")
		filePath = "../input/input.example.txt"
	}

	day18(filePath, 1, 2)
}

func day18(filePath string, parts ...int) {
	board := readInput(filePath)

	for _, part := range parts {
		if part == 1 {
			part1(board)
		} else if part == 2 {
			part2(board)
		}
	}
}

func findShortestPath(board Board) Node {
	startingNode := Node{pos: Vector{X: 0, Y: 0}, parent: nil}
	queue := []Node{startingNode}
	visited := make(map[Vector]bool)
	target := Vector{X: board.bounds.X - 1, Y: board.bounds.Y - 1}

	var shortestPath Node
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.pos] {
			continue
		}
		if current.pos == target {
			shortestPath = current
			break
		}

		for _, neighbor := range current.pos.Neighbors() {
			if neighbor.InBounds(board.bounds.X, board.bounds.Y) && board.board[neighbor.Y][neighbor.X] != "#" && !visited[neighbor] {
				queue = append(queue, Node{pos: neighbor, parent: &current})
			}
		}

		visited[current.pos] = true
	}

	return shortestPath
}

func part1(board Board) {
	// Generate the falling bytes
	board.generateFallingBytes(1024)

	// Set up the TUI program
	p := tea.NewProgram(tui.NewModel("Day 18 - Part One"), tea.WithAltScreen(), tea.WithMouseCellMotion())
	grid := plotBoard(board)

	go func() {
		// Display the initial board
		p.Send(tui.UpdateViewport("Finding the shortest path through the grid: \n\n"+grid, 1000))
		time.Sleep(100 * time.Millisecond)

		// Find the shortest path using BFS
		shortestPath := findShortestPath(board)

		// Mark the shortest path on the board for visualization. (0's mark the visited blocks)
		length := 0
		for shortestPath.parent != nil {
			board.board[shortestPath.pos.Y][shortestPath.pos.X] = tui.RobotStyle.Render("0")
			shortestPath = *shortestPath.parent
			length++
		}
		board.board[shortestPath.pos.Y][shortestPath.pos.X] = tui.RobotStyle.Render("0") // Mark the start line as well.

		// Formatting the tui output.
		s := fmt.Sprintf("Shortest Path Length: %d\n", length)
		s += plotBoard(board)
		s += "\nPress 'q' or 'ctrl+c' or 'esc' to exit (will move on to part 2)...."

		p.Send(tui.UpdateViewport(s, 2*len(board.board[0]))) // Update the terminal ui.
	}()

	// Run the terminal UI. has to be in a separate thread.
	if _, err := p.Run(); err != nil {
		fmt.Println("error running visualizer: ", err)
		os.Exit(1)
	}
}

func part2(board Board) {
	// Need to find the first falling byte that blocks the shortest path.
	p := tea.NewProgram(tui.NewModel("Day 18 - Part Two"), tea.WithAltScreen(), tea.WithMouseCellMotion())
	solution := Vector{X: -1, Y: -1}
	go func() {
		// Copy the board
		fallingBytes := board.fallingBytes
		for i := 0; i <= len(fallingBytes); i++ {
			board.generateFallingBytes(1)
			board.fallingBytes = board.fallingBytes[1:]

			// Build the viewport update
			viewport := "Attempting to find the first falling byte that blocks all solutions\n"
			viewport += fmt.Sprintf(" Current Falling Byte: (%d,/%d) ................\n", i, len(fallingBytes))
			viewport += "  (Might look stuck when it's near the answer)\n\n"
			viewport += plotBoard(board)

			// Send the update to the viewport.
			p.Send(tui.UpdateViewport(viewport, 1000))
			// time.Sleep(1 * time.Millisecond) // Sleeping for 1 millisec for brevity.

			// Check if the shortest path is blocked.
			shortestPath := findShortestPath(board)
			target := Vector{X: board.bounds.X - 1, Y: board.bounds.Y - 1}
			if shortestPath.pos != target {
				p.Send(tui.UpdateViewport(fmt.Sprintf("Sol: %d,%d\n%s", fallingBytes[i].X, fallingBytes[i].Y, plotBoard(board)), 1000))
				log.Println(fallingBytes[i])
				solution = fallingBytes[i]
				break
			}
		}

		if solution.X == -1 {
			p.Send(tui.UpdateViewport("No solution", 1000))
			log.Fatal("no solution")
		}

		// Solution to the puzzle exists. Let's draw out the shortest path before that byte fell, draw the path, and mark the final byte
		board.board[solution.Y][solution.X] = "." // remove that byte first
		shortestPath := findShortestPath(board)   // find the shortest path
		// Mark the path
		for shortestPath.parent != nil {
			board.board[shortestPath.pos.Y][shortestPath.pos.X] = tui.RobotStyle.Render("0")
			shortestPath = *shortestPath.parent
		}
		board.board[shortestPath.pos.Y][shortestPath.pos.X] = tui.RobotStyle.Render("0") // Mark the start line as well.

		// Mark the last fallen byte.
		board.board[solution.Y][solution.X] = "@"

		// Draw the board.
		s := fmt.Sprintf("Solution: %d,%d\n\n%s", solution.X, solution.Y, plotBoard(board)) // construct the viewport
		p.Send(tui.UpdateViewport(s, 1000))                                                 //draw the viewport

	}()

	if _, err := p.Run(); err != nil {
		log.Fatal("part two: error with TUI", err)
	}
}

func (board *Board) generateFallingBytes(amountToGenerate int) {
	if amountToGenerate > len(board.fallingBytes) {
		return
	}

	// Update the board and generate the bytes on the grid.
	if len(board.board) == 0 {
		board.board = make([][]string, 0, board.bounds.Y)
		for range board.bounds.Y {
			board.board = append(board.board, make([]string, board.bounds.X))
		}
	}

	for i := 0; i < amountToGenerate; i++ {
		vec := board.fallingBytes[i]
		board.board[vec.Y][vec.X] = "#"
	}
}

func plotBoard(board Board) string {
	// Plot the board using the tui.
	out := ""
	for _, row := range board.board {
		for _, cell := range row {
			if cell == "" {
				out += "."
			} else if cell == "#" {
				out += tui.WallStyle.Render("#")
			} else {
				out += cell
			}
		}
		out += "\n"
	}

	return out
}

func readInput(filePath string) Board {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bounds := readBounds(scanner)

	// Skip empty line
	scanner.Scan()

	fallingBytes := readFallingBytes(scanner)

	return Board{
		bounds:       bounds,
		fallingBytes: fallingBytes,
	}
}

func readBounds(scanner *bufio.Scanner) Vector {
	scanner.Scan()
	rawBounds := strings.Split(scanner.Text(), ",")

	boundX, err := strconv.Atoi(rawBounds[0])
	if err != nil {
		log.Fatal("error parsing bounds(x): ", err)
	}

	boundY, err := strconv.Atoi(rawBounds[1])
	if err != nil {
		log.Fatal("error parsing bounds(y): ", err)
	}

	bounds := Vector{X: boundX + 1, Y: boundY + 1}
	return bounds
}

func readFallingBytes(scanner *bufio.Scanner) []Vector {
	var fallingBytes []Vector

	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ",")
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			log.Fatal("error parsing falling byte(x): ", err)
		}

		y, err := strconv.Atoi(coords[1])
		if err != nil {
			log.Fatal("error parsing falling byte(y): ", err)
		}

		fallingBytes = append(fallingBytes, Vector{X: x, Y: y})
	}
	return fallingBytes
}
