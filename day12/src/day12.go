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
	log.SetFlags(0)

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

	buffer := make([][]string, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		buffer = append(buffer, line)
	}
	fmt.Println("Original Input: ", buffer)
	coordinateRegions := make(map[Vector]int)
	plotGarden(buffer, coordinateRegions)

	// AREA - seems to just be the count of garden plots (letters in contiguous space)
	// PERIMETER - more annoying to do. is the actual perimeter. 2*(x2-x1 + y2-y1)
	regions := make([][]Vector, 0)
	visited := make(map[Vector]bool)
	for i, row := range buffer {
		for j, plot := range row {
			coordinate := Vector{x: j, y: i}

			// I need to find the coordinates for that region. assign them to that region.
			// If the current coordinate has a region, not do it at all
			if visited[coordinate] {
				continue
			}

			// DFS to find neighboring plots of this new region.
			// visited := make(map[Vector]bool)
			region := make([]Vector, 0, 1)
			stack := []Vector{coordinate}

			for len(stack) > 0 {
				pos := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				// if visited, it already has a region. continue.
				if visited[pos] {
					continue
				}

				// coordinateRegions[pos] = plot
				visited[pos] = true
				region = append(region, pos)

				for _, neighbor := range pos.neighbors() {
					if neighbor.inBounds(len(row), len(buffer)) && !visited[neighbor] && buffer[neighbor.y][neighbor.x] == plot {
						stack = append(stack, neighbor)
					}
				}
			}

			regions = append(regions, region)
			for _, pos := range region {
				coordinateRegions[pos] = len(regions)
			}
			// plotGarden(buffer, coordinateRegions)
		}
	}
	fmt.Println(regions)

	// Part one: PRICE
	price := 0
	for _, region := range regions {
		area := len(region)
		per := perimeter(coordinateRegions, region)

		price += area * per
		// log.Println("REGION: ", buffer[region[0].y][region[0].x], " Perimeter: ", per, "Area: ", area, " Price: ")
	}
	fmt.Println("PART-ONE: ", price)

	// Part Two : Discount on perimeter
	price = 0
	for _, region := range regions {
		area := len(region)
		per := perimeter2(coordinateRegions, region)

		price += area * per
		log.Println("REGION: ", buffer[region[0].y][region[0].x], " Perimeter: ", per, "Area: ", area, " Price: ")
	}
	fmt.Println("PART-TWO: ", price)

}

func perimeter(regionCoordinates map[Vector]int, region []Vector) int {
	regionID := regionCoordinates[region[0]]

	perimeter := 0
	for _, pos := range region {
		for _, neighbor := range pos.neighbors() {
			if regionCoordinates[neighbor] != regionID {
				perimeter++
			}
		}
	}
	return perimeter
}

func perimeter2(regionCoordinates map[Vector]int, region []Vector) int {
	// This one is inspired/copied/cheated from reddit. The idea is at least.
	// To find the "perimeter", which in this case is the number of sides(part 2),
	// just count corners. Count each corner once, that's the perimeter.

	regionID := regionCoordinates[region[0]]
	perimeter := 0
	for _, pos := range region {
		// neighbor above
		neighbor := pos.up()
		// 4 different corners to check, each with their own conditions.
		if regionCoordinates[neighbor] != regionID { // First check if it's in the region.
			if regionCoordinates[pos.left()] != regionID { // |`` (external border)
				perimeter++
			}
			if regionCoordinates[pos.right()] != regionID { // ``| (external border)
				perimeter++
			}
			if regionCoordinates[pos.right()] == regionID && regionCoordinates[neighbor.right()] == regionID { // |__	(internal border)
				perimeter++
			}
			if regionCoordinates[pos.left()] == regionID && regionCoordinates[neighbor.left()] == regionID { // __| (internal border)
				perimeter++
			}
		}

		neighbor = pos.down()
		// 4 different corners to check, each with their own conditions.
		if regionCoordinates[neighbor] != regionID { // First check if it's in the region.
			if regionCoordinates[pos.left()] != regionID { //
				perimeter++
			}
			if regionCoordinates[pos.right()] != regionID { // top-right corner
				perimeter++
			}
			if regionCoordinates[pos.right()] == regionID && regionCoordinates[neighbor.right()] == regionID { // shape of bottom-right corner ( |_ )
				perimeter++
			}
			if regionCoordinates[pos.left()] == regionID && regionCoordinates[neighbor.left()] == regionID {
				perimeter++
			}
		}
	}
	return perimeter
}

func (vec Vector) neighbors() []Vector {
	return []Vector{
		vec.up(),
		vec.left(),
		vec.down(),
		vec.right(),
	}
}

func (vec Vector) down() Vector {
	vec.y += 1
	return vec
}
func (vec Vector) up() Vector {
	vec.y -= 1
	return vec
}
func (vec Vector) right() Vector {
	vec.x += 1
	return vec
}
func (vec Vector) left() Vector {
	vec.x -= 1
	return vec
}
func (vec Vector) inBounds(xBounds int, yBounds int) bool {
	if vec.x < 0 || vec.y < 0 || vec.x >= xBounds || vec.y >= yBounds {
		return false
	}
	return true
}

func plotGarden(gardenMap [][]string, regionCoordinates map[Vector]int) {
	gardenViz := make([]string, 0)
	gardenViz = append(gardenViz, "+", strings.Repeat("-", len(gardenMap[0])*2)+" +")
	for _, row := range gardenMap {
		out := "|  "
		for _, plot := range row {
			out += plot + " "
		}
		gardenViz = append(gardenViz, out+"|")
		gardenViz = append(gardenViz, "|"+strings.Repeat(" ", len(gardenMap[0])*2+2)+"|")
	}
	gardenViz = append(gardenViz, "+", strings.Repeat("-", len(gardenMap[0])*2)+" +")

	for _, line := range gardenViz {
		log.Println(line)
	}
}
