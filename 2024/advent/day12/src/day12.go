package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func day12(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}

	scanner := bufio.NewScanner(file)
	buffer := make([][]string, 0, 1024)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		buffer = append(buffer, line)
	}

	regions, coordinateRegions := findRegions(buffer)
	plotGarden(buffer)

	// Part one: PRICE
	price := price_partOne(regions, coordinateRegions)
	fmt.Println("PART-ONE: ", price)

	// Part Two : Discount on perimeter
	price = price_partTwo(regions, coordinateRegions)
	fmt.Println("PART-TWO: ", price)
}

func price_partOne(regions [][]Vector, coordRegions map[Vector]int) int {
	price := 0
	for _, region := range regions {
		area := len(region)
		per := perimeter(coordRegions, region)

		price += area * per
	}

	return price
}

func price_partTwo(regions [][]Vector, coordRegions map[Vector]int) int {
	price := 0
	for _, region := range regions {
		area := len(region)
		per := perimeter2(coordRegions, region)

		price += area * per
	}

	return price
}

func findRegions(board [][]string) (regions [][]Vector, coordRegions map[Vector]int) {
	regions = make([][]Vector, 0)
	coordRegions = make(map[Vector]int)

	visited := make(map[Vector]bool)
	for i, row := range board {
		for j, plot := range row {
			coordinate := Vector{x: j, y: i}
			if visited[coordinate] {
				continue
			}
			region := make([]Vector, 0, 1)
			stack := []Vector{coordinate}

			for len(stack) > 0 {
				pos := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if visited[pos] {
					continue
				}
				visited[pos] = true
				region = append(region, pos)

				for _, neighbor := range pos.neighbors() {
					if neighbor.inBounds(len(row), len(board)) && !visited[neighbor] && board[neighbor.y][neighbor.x] == plot {
						stack = append(stack, neighbor)
					}
				}
			}

			regions = append(regions, region)
			for _, pos := range region {
				coordRegions[pos] = len(regions)
			}
			// plotGarden(buffer, coordinateRegions)
		}
	}

	return regions, coordRegions
}

func plotGarden(gardenMap [][]string) {
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
