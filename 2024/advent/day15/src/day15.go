package main

// Too many imports, ik.
import (
	"advent2024/tui"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var robotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))
var wallStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
var boxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0AA00"))

type Vector struct {
	x int
	y int
}

var (
	box   = "O"
	robot = "@"
	empty = "."
	wall  = "#"

	box2_1 = "["
	box2_2 = "]"
	robot2 = "@."
	empty2 = ".."
	wall2  = "##"
)
var sth = time.Duration(100)

type Board struct {
	grid          [][]string
	robotPosition Vector
	xBounds       int
	yBounds       int
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

		// log.SetOutput(logFile)
		log.SetFlags(0)
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
	grid, robotPos, robotMoves := readInput(filePath)

	log.Println("Original Input: ")
	log.Println(plotMap(grid))

	log.Println("Robot Moves:")
	log.Println(robotMoves)

	// Part One
	// partOne(grid, robotMoves, robotPos)

	partTwo(grid, robotMoves, robotPos)
}
func partOne(grid [][]string, robotMoves []string, robotPos Vector) {
	// Needs to be able to shove boxes when there is no wall behind them
	board := Board{
		robotPosition: robotPos,
		grid:          grid,
		xBounds:       len(grid[0]),
		yBounds:       len(grid),
	}

	p := tea.NewProgram(tui.Model{}, tea.WithAltScreen(), tea.WithMouseCellMotion())

	gpsScore := 0
	go func() {
		// Let's track the robot's position
		for _, move := range robotMoves {
			board.moveRobot(move)
			gpsScore = board.gps()
			// s := fmt.Sprintf("Move %s: %d/%d \nRobot Position: (%d,%d)\n%s\n\ngps: %d", move, i+1, len(robotMoves), board.robotPosition.x, board.robotPosition.y, plotMap(board.grid), gpsScore)
			// p.Send(tui.UpdateViewport(s, 100))
			// time.Sleep(0 * time.Millisecond)
		}

		s := fmt.Sprintf("%s\ngps: %d\n", plotMap(board.grid), gpsScore)
		p.Send(tui.UpdateViewport(s, 1000))
		fmt.Println(gpsScore)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("error running visualizer: ", err)
		os.Exit(1)
	}
}

func partTwo(grid [][]string, robotMoves []string, robotPos Vector) {
	// Needs to be able to shove boxes when there is no wall behind them
	updatedGrid, robotPos := transformGrid(grid)
	fmt.Println(plotMap(updatedGrid))

	board := Board{
		robotPosition: robotPos,
		grid:          updatedGrid,
		xBounds:       len(updatedGrid[0]),
		yBounds:       len(updatedGrid),
	}

	p := tea.NewProgram(tui.Model{}, tea.WithAltScreen(), tea.WithMouseCellMotion())

	gpsScore := 0
	go func() {
		// Let's track the robot's position
		for i, move := range robotMoves {
			s := fmt.Sprintf("Move %s: %d/%d \nRobot Position: (%d,%d)\n%s\n\ngps: %d", move, i+1, len(robotMoves), board.robotPosition.x, board.robotPosition.y, plotMap(board.grid), gpsScore)
			gpsScore = board.gps()
			p.Send(tui.UpdateViewport(s, 100))
			time.Sleep(100 * time.Millisecond)
			board.moveRobot(move)
			gpsScore = board.gps()
			s = fmt.Sprintf("Move %s: %d/%d \nRobot Position: (%d,%d)\n%s\n\ngps: %d", move, i+1, len(robotMoves), board.robotPosition.x, board.robotPosition.y, plotMap(board.grid), gpsScore)
			p.Send(tui.UpdateViewport(s, 100))
			time.Sleep(10 * time.Millisecond)
		}
		// gpsScore = board.gps()
		s := fmt.Sprintf("%s\ngps: %d\n", plotMap(board.grid), gpsScore)
		p.Send(tui.UpdateViewport(s, 1000))
		fmt.Println(gpsScore)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("error running visualizer: ", err)
		os.Exit(1)
	}
}

func readInput(filePath string) (grid [][]string, robotPosition Vector, robotMoves []string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	// Input file has two sections: A map and a sequence of robot moves

	// Read the map
	buffer := make([][]string, 0, 1024)
	robotFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		splitLine := strings.Split(line, "") // Separate the characters.

		buffer = append(buffer, splitLine)
		// Find the robot position, if not already tracked
		if robotFound { // this just avoids further iteration thru the lines if the robot is already found.
			continue
		}
		for i, char := range splitLine {
			if char == "@" { // Robot Found
				robotPosition.y = len(buffer) - 1
				robotPosition.x = i
				robotFound = true
			}
		}
	}

	// Read the robot's movements
	rawRobotMoves := ""
	for scanner.Scan() {
		rawRobotMoves += scanner.Text()
	}

	robotMoves = strings.Split(rawRobotMoves, "")
	return buffer, robotPosition, robotMoves
}

func transformGrid(grid [][]string) ([][]string, Vector) {
	updatedGrid := make([][]string, 0, len(grid))
	robotPos := Vector{}

	for j, row := range grid {
		newRow := make([]string, 0)
		for _, char := range row {
			if char == "O" {
				newRow = append(newRow, "[", "]")
			} else if char == "@" {
				robotPos.x = len(newRow)
				robotPos.y = j
				newRow = append(newRow, "@", ".")
			} else {
				newRow = append(newRow, char, char)
			}
		}
		updatedGrid = append(updatedGrid, newRow)
	}
	return updatedGrid, robotPos
}

func (board *Board) moveRobot(move string) {
	// Needs to be able to shove boxes when there is no wall behind them
	// Let's track the robot's position

	targetCoord := board.robotPosition.add(getMove(move))
	if !targetCoord.inBounds(board.xBounds, board.yBounds) {
		return
	}
	target := board.grid[targetCoord.y][targetCoord.x]

	switch target {
	case box:
		if board.moveBox(move, targetCoord) {
			board.grid[targetCoord.y][targetCoord.x] = robot
			board.grid[board.robotPosition.y][board.robotPosition.x] = empty
		} else {
			return
		}
	case box2_1:
		temp := board.copy()
		if board.moveBox(move, targetCoord.add(Vector{x: 1})) && board.moveBox(move, targetCoord) {
			board.grid[targetCoord.y][targetCoord.x] = robot
			board.grid[board.robotPosition.y][board.robotPosition.x] = empty
		} else {
			board.grid = temp.grid
			board.robotPosition = temp.robotPosition
			return
		}

	case box2_2:
		temp := board.copy()
		if board.moveBox(move, targetCoord.add(Vector{x: -1})) && board.moveBox(move, targetCoord) {
			board.grid[targetCoord.y][targetCoord.x] = robot
			board.grid[board.robotPosition.y][board.robotPosition.x] = empty
		} else {
			board.grid = temp.grid
			board.robotPosition = temp.robotPosition
			return
		}
	case wall:
		return
	case empty:
		board.grid[targetCoord.y][targetCoord.x] = robot
		board.grid[board.robotPosition.y][board.robotPosition.x] = empty
	}

	board.robotPosition = targetCoord
}

func (board *Board) moveBox(move string, startPos Vector) bool {
	targetCoord := startPos.add(getMove(move))
	if !targetCoord.inBounds(board.xBounds, board.yBounds) {
		return false
	}
	target := board.grid[targetCoord.y][targetCoord.x]

	switch target {
	case box:
		if board.moveBox(move, targetCoord) {
			board.grid[targetCoord.y][targetCoord.x] = box
			board.grid[startPos.y][startPos.x] = empty
			return true
		}
	case box2_1:
		temp := board.copy()
		if board.moveBox(move, targetCoord.add(Vector{x: 1})) && board.moveBox(move, targetCoord) {
			board.grid[targetCoord.y][targetCoord.x] = board.grid[startPos.y][startPos.x]
			board.grid[startPos.y][startPos.x] = empty
			return true
		} else {
			board.grid = temp.grid
			board.robotPosition = temp.robotPosition
			return false
		}

	case box2_2:
		temp := board.copy()
		if board.moveBox(move, targetCoord.add(Vector{x: -1})) && board.moveBox(move, targetCoord) {
			board.grid[targetCoord.y][targetCoord.x] = board.grid[startPos.y][startPos.x]
			board.grid[startPos.y][startPos.x] = empty
			return true
		} else {
			board.grid = temp.grid
			board.robotPosition = temp.robotPosition
			return false
		}

	case wall:
		return false
	case empty:
		board.grid[targetCoord.y][targetCoord.x] = board.grid[startPos.y][startPos.x]
		board.grid[startPos.y][startPos.x] = empty
		return true
	}
	return false
}

func (board Board) copy() Board {
	newGrid := make([][]string, len(board.grid))
	for i, r := range board.grid {
		newGrid[i] = make([]string, 0, len(r))
		newGrid[i] = append(newGrid[i], r...)
	}

	newPos := Vector{x: board.robotPosition.x, y: board.robotPosition.y}
	return Board{
		grid:          newGrid,
		robotPosition: newPos,
		xBounds:       board.xBounds,
		yBounds:       board.yBounds,
	}
}

func (board Board) gps() int {
	score := 0
	for j, row := range board.grid {
		for i, char := range row {
			if char == box || char == "[" {
				score += (100 * j) + i
			}
		}
	}

	return score
}

func (vec Vector) add(vec2 Vector) Vector {
	vec.x += vec2.x
	vec.y += vec2.y
	return vec
}

func (vec Vector) inBounds(xBounds int, yBounds int) bool {
	if vec.x < 0 || vec.x > xBounds || vec.y < 0 || vec.y > yBounds {
		return false
	}
	return true
}

func getMove(move string) Vector {
	robotMove := Vector{}
	switch move {
	case "^":
		robotMove.y = -1
	case "v":
		robotMove.y = 1
	case ">":
		robotMove.x = 1
	case "<":
		robotMove.x = -1
	}

	return robotMove
}

func plotMap(givenMap [][]string) string {
	s := ""
	for _, row := range givenMap {
		for _, c := range row {
			switch c {
			case wall:
				s += wallStyle.Render(wall) + " "
			case robot:
				s += robotStyle.Render(robot) + " "
			case box:
				s += boxStyle.Render(box) + " "
			default:
				s += c + " "
			}
		}
		s += "\n"
	}
	return s
}
