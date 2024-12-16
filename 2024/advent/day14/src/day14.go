package main

// Too many imports, ik.
import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

type Vector struct {
	x int
	y int
}

type Robot struct {
	pos Vector
	vel Vector
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

	mapWidth := 11
	mapHeight := 7

	var filePath string
	if !runExampleInput {
		fmt.Println("Running Test input...")
		filePath = "../input/input.txt"
		mapWidth = 101
		mapHeight = 103
	} else {
		fmt.Println("Running Example input...")
		filePath = "../input/input.example.txt"
	}

	mapBounds := Vector{
		x: mapWidth,
		y: mapHeight,
	}

	day12(filePath, mapBounds)
}

func day12(filePath string, mapBounds Vector) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}

	scanner := bufio.NewScanner(file)
	robots := make([]Robot, 0)
	for scanner.Scan() {
		line := scanner.Text()
		robots = append(robots, makeRobot(line))
	}
	log.Println(len(robots), robots)

	// // Part One
	updatedRobots := moveRobots(robots, mapBounds, 100)
	log.Println(len(updatedRobots), updatedRobots)

	log.Println("Original Map")
	plotMap(robots, mapBounds)

	log.Println("UPDATED MAP")
	// plotMap(updatedRobots, mapBounds)

	fmt.Println(safetyScore(updatedRobots, mapBounds))

	// Part Two
	// christmasTree(robots, mapBounds)

	// since i know my answer, just to see the log
	plotMap(moveRobots(robots, mapBounds, 7858), mapBounds)
}

func christmasTree(robots []Robot, mapBounds Vector) {
	// I HAVE NO FUCKING CLUE HOW TO DO THIS PROGRAMATICALLY
	i := 1
	for i <= 10000 {
		// Move the robots by one second
		robots = moveRobots(robots, mapBounds, 1)

		// Generate an image for the new robot positions
		plotMap2(robots, mapBounds, i, "./out")
		i++
	}
}

func moveRobots(robots []Robot, mapBounds Vector, timeSeconds int) []Robot {
	updatedRobots := make([]Robot, 0, len(robots))
	for _, robot := range robots {
		updatedRobots = append(updatedRobots, robot.moveRobot(timeSeconds, mapBounds))
	}
	return updatedRobots
}

func plotMap2(robots []Robot, mapBounds Vector, timePassed int, dir string) {
	fieldMap := make([][]string, mapBounds.y)
	positions := make(map[Vector]int)

	for _, robot := range robots {
		positions[robot.pos]++
	}

	// let's make an image
	img := image.NewRGBA(image.Rect(0, 0, mapBounds.x, mapBounds.y))
	pixelColor := color.RGBA{255, 255, 255, 255}
	bgColor := color.RGBA{0, 0, 0, 255}

	for i, row := range fieldMap {
		row = make([]string, mapBounds.x)
		for j, _ := range row {
			coord := Vector{x: j, y: i}
			if positions[coord] != 0 {
				row[j] = strconv.Itoa(positions[coord])
				img.Set(j, i, pixelColor)
			} else {
				row[j] = " "
				img.Set(j, i, bgColor)
			}
		}
		fieldMap[i] = row

		rowStr := ""
		for _, r := range row {
			rowStr += r
		}
		// log.Println(row)
	}
	file, err := os.Create(fmt.Sprintf("%s/%d.png", dir, timePassed))
	if err != nil {
		log.Fatal("error creating file for img")
	}
	defer file.Close()
	png.Encode(file, img)

}
func plotMap(robots []Robot, mapBounds Vector) {
	fieldMap := make([][]string, mapBounds.y)
	positions := make(map[Vector]int)

	for _, robot := range robots {
		positions[robot.pos]++
	}

	for i, row := range fieldMap {
		row = make([]string, mapBounds.x)
		for j, _ := range row {
			coord := Vector{x: j, y: i}
			if positions[coord] != 0 {
				// row[j] = strconv.Itoa(positions[coord])
				row[j] = "@"
			} else {
				row[j] = " "
			}
		}
		fieldMap[i] = row

		rowStr := ""
		for _, r := range row {
			rowStr += r
		}
		log.Println(rowStr)
	}
	log.Println()
}

func safetyScore(robots []Robot, mapBounds Vector) int {
	middleV := (mapBounds.x - 1) / 2
	middleH := (mapBounds.y - 1) / 2

	quadrant1 := 0
	quadrant2 := 0
	quadrant3 := 0
	quadrant4 := 0

	for _, robot := range robots {
		x := robot.pos.x
		y := robot.pos.y

		// Quadrant 1
		if x < middleV && y < middleH {
			quadrant1++
		}
		// Quadrant 2
		if x > middleV && y < middleH {
			quadrant2++
		}
		// Quadrant 3
		if x < middleV && y > middleH {
			quadrant3++
		}
		// Quadrant 4
		if x > middleV && y > middleH {
			quadrant4++
		}
	}

	fmt.Println(quadrant1)
	fmt.Println(quadrant2)
	fmt.Println(quadrant3)
	fmt.Println(quadrant4)

	return quadrant1 * quadrant2 * quadrant3 * quadrant4
}

func (robot Robot) moveRobot(timeSeconds int, mapBounds Vector) Robot {
	robot.pos.x = (robot.pos.x + (robot.vel.x * timeSeconds))
	robot.pos.y = (robot.pos.y + (robot.vel.y * timeSeconds))

	if robot.pos.x < 0 {
		robot.pos.x = (mapBounds.x - (-robot.pos.x % mapBounds.x)) % mapBounds.x
	} else if robot.pos.x > mapBounds.x {
		robot.pos.x %= mapBounds.x
	}

	if robot.pos.y < 0 {
		robot.pos.y = (mapBounds.y - (-robot.pos.y % mapBounds.y)) % mapBounds.y
	} else if robot.pos.y > mapBounds.y {
		robot.pos.y %= mapBounds.y
	}

	return robot
}

func makeRobot(str string) Robot {
	temp := strings.Split(str, " ")
	if len(temp) != 2 {
		log.Fatal("error parsing robot params: expecting p=x,y v=x,y")
	}

	robot := Robot{}

	// Parse Position
	rawPos := strings.Split(strings.Split(temp[0], "=")[1], ",")
	if len(temp) != 2 {
		log.Fatal("error parsing robot pos: expecting p=x,y")
	}

	xPos, err := strconv.Atoi(rawPos[0])
	if err != nil {
		log.Fatal("error parsing robot X pos: NAN")
	}
	yPos, err := strconv.Atoi(rawPos[1])
	if err != nil {
		log.Fatal("error parsing robot Y pos: NAN")
	}
	robot.pos.x = xPos
	robot.pos.y = yPos

	// Parse Velocity
	rawVel := strings.Split(strings.Split(temp[1], "=")[1], ",")
	if len(temp) != 2 {
		log.Fatal("error parsing robot vel: expecting v=x,y")
	}

	xVel, err := strconv.Atoi(rawVel[0])
	if err != nil {
		log.Fatal("error parsing robot X vel: NAN")
	}
	yVel, err := strconv.Atoi(rawVel[1])
	if err != nil {
		log.Fatal("error parsing robot Y vel: NAN")
	}
	robot.vel.x = xVel
	robot.vel.y = yVel

	return robot
}
