package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	runExampleFlag := flag.Bool("e", true, "specifies which input file to use: example input or test input")
	flag.Parse()

	runExample := *runExampleFlag

	filePath := "../input/input.example.txt"
	if !runExample {
		filePath = "../input/input.txt"
	}

	day8(filePath)
}

type Vector struct {
	x int
	y int
}

func day8(filePath string) {
	// Setup logging.
	logFile, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Open the input file.
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Input file not found: ", err)
	}
	scanner := bufio.NewScanner(file)

	// Scan the input map.
	inputMap := make([][]string, 0, 1024)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, "")
		inputMap = append(inputMap, splitLine)
	}

	// Logging the initial Map.
	log.Println("-- Original Map")
	logMap(inputMap)

	part_one(inputMap)
	part_two(inputMap)
}

func logMap(givenMap [][]string) {
	for _, l := range givenMap {
		log.Println(l)
	}
}

func part_one(inputMap [][]string) {
	// Part One
	// Let's map out and grab the locations of the antennae
	locations := make(map[string][]Vector)
	for i, row := range inputMap {
		for j, col := range row {
			if col == "." {
				continue
			}

			store, exists := locations[col]
			if !exists {
				locations[col] = make([]Vector, 0, 1)
				store = locations[col]
			}
			store = append(store, Vector{x: j, y: i})
			locations[col] = store
		}
	}
	log.Println(locations)

	log.Println("---------------------------------------------------------------")
	log.Println("Drawing Antinodes")
	mapYBound := len(inputMap)
	mapXBound := len(inputMap[0])

	// Let's make the anti-nodes
	locations["#"] = make([]Vector, 0)
	antinodes := locations["#"]
	for _, vals := range locations {
		// for each antenna "type", i have a list of vectors for the antenna positions.
		// I'm going to loop over that, for each element i'm looping through, make antinodes for all the others that come after it.
		// nlogn? not quite n^2, so it should be nlogn
		if len(vals) < 2 {
			continue
		}
		for i, antenna1 := range vals {
			for _, antenna2 := range vals[i+1:] {
				distance := Vector{
					x: antenna1.x - antenna2.x,
					y: antenna1.y - antenna2.y,
				}

				// First antinode
				anti1 := Vector{ // antinode 1: for antenna 1
					x: antenna1.x + distance.x,
					y: antenna1.y + distance.y,
				}
				anti2 := Vector{ // antinode 2: for antenna 2
					x: antenna2.x - distance.x,
					y: antenna2.y - distance.y,
				}

				// Check if the antinodes are valid (on the map)
				if anti1.x >= 0 && anti1.x < mapXBound && anti1.y >= 0 && anti1.y < mapYBound {
					antinodes = append(antinodes, anti1)
					inputMap[anti1.y][anti1.x] = "#"
				}
				if anti2.x >= 0 && anti2.x < mapXBound && anti2.y >= 0 && anti2.y < mapYBound {
					antinodes = append(antinodes, anti2)
					inputMap[anti2.y][anti2.x] = "#"
				}

				logMap(inputMap)
				log.Println()
				log.Println("----------------------------------------------------------")
			}
		}
	}
	log.Println(antinodes)
	locations["#"] = antinodes

	// Let's count the unique antinodes for the final answer. Again, go doesn't have a set :(
	uniqueAntinodes := make(map[Vector]bool)
	for _, antinode := range antinodes {
		uniqueAntinodes[antinode] = true
	}

	fmt.Println("PART ONE SOLUTION: ", len(uniqueAntinodes))
	log.Println("PART ONE SOLUTION: ", len(uniqueAntinodes))
}

func part_two(inputMap [][]string) {
	// Let's map out and grab the locations of the antennae
	locations := make(map[string][]Vector)
	for i, row := range inputMap {
		for j, col := range row {
			if col == "." {
				continue
			}

			store, exists := locations[col]
			if !exists {
				locations[col] = make([]Vector, 0, 1)
				store = locations[col]
			}
			store = append(store, Vector{x: j, y: i})
			locations[col] = store
		}
	}
	log.Println(locations)

	log.Println("-----------------------PART TWO---------------------------------")
	log.Println("Drawing Antinodes")
	mapYBound := len(inputMap)
	mapXBound := len(inputMap[0])

	// Let's make the anti-nodes
	locations["#"] = make([]Vector, 0)
	antinodes := locations["#"]
	for _, vals := range locations {
		if len(vals) < 2 {
			continue
		}
		for i, antenna1 := range vals {
			for _, antenna2 := range vals[i+1:] {
				distance := Vector{
					x: antenna1.x - antenna2.x,
					y: antenna1.y - antenna2.y,
				}

				// This time, we don't just need 2 antinodes, but as many as possible within the bounds of the map
				// First Antenna antinodes
				for multiplier := 0; ; multiplier++ {
					anti1 := Vector{ // antinode 1: for antenna 1
						x: antenna1.x + multiplier*(distance.x),
						y: antenna1.y + multiplier*(distance.y),
					}

					// Check if the antinode is valid (on the map)
					if anti1.x >= 0 && anti1.x < mapXBound && anti1.y >= 0 && anti1.y < mapYBound {
						antinodes = append(antinodes, anti1)
						inputMap[anti1.y][anti1.x] = "#"
					} else {
						// nothing after this will be in bounds after either.
						break
					}
				}

				// Second Antenna Antinodes
				for multiplier := 0; ; multiplier++ {
					anti2 := Vector{ // antinode 2: for antenna 2
						x: antenna2.x - multiplier*(distance.x),
						y: antenna2.y - multiplier*(distance.y),
					}

					// Check if the antinode is valid (on the map)
					if anti2.x >= 0 && anti2.x < mapXBound && anti2.y >= 0 && anti2.y < mapYBound {
						antinodes = append(antinodes, anti2)
						inputMap[anti2.y][anti2.x] = "#"
					} else {
						break
					}
				}

				logMap(inputMap)
				log.Println()
				log.Println("----------------------------------------------------------")
			}
		}
	}
	log.Println(antinodes)
	locations["#"] = antinodes

	// Let's count the unique antinodes for the final answer. Again, go doesn't have a set :(
	uniqueAntinodes := make(map[Vector]bool)
	for _, antinode := range antinodes {
		uniqueAntinodes[antinode] = true
	}

	fmt.Println("PART TWO SOLUTION: ", len(uniqueAntinodes))
	log.Println("PART TWO SOLUTION: ", len(uniqueAntinodes))
}
