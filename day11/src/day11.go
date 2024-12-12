package main

// Too many imports, ik.
import (
	"os"
	"fmt"
	"log"
	"flag"
	"math"
	"bufio"
	"strconv"
	"strings"
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

	day11(filePath)
}

func day11(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([]int, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		for _, v := range line {
			val, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("error parsing int: %v", err)
			}

			buffer = append(buffer, val)
		}
	}
	log.Println("Original Input: ", buffer)

	// BLINK!!!!!!!
	stones := make([]int, 0, 1024)
	stones = append(stones, buffer...)

	mapStones := make(map[int]int)
	for _, stone := range buffer {
		mapStones[stone]++
	}

	for range 75 {
		// blink(&stones)
		mapStones = blink2(mapStones)
	}

	fmt.Println(len(stones))

	count := 0
	for _, stoneCount := range mapStones {
		count += stoneCount
	}
	fmt.Println(count)

}

func blink2(stones map[int]int) map[int]int {
	/*
		PART-TWO Solution
		Couldn't use the arrays from part one, they got too large and caused memory crashes.
		This time, stores the stones in a map, mapping the stones number => number of stones with that number.
		This will reduce storage, as a stone with a count of thousands will need just a single field, it's ID and count, instead of #count array slots.
		It will also allow for bulk updating, for each stone value, we can update all the counts instead of one at a time
	*/

	// Need a snapshot of the array at the start of the blink. This will prevent any updates to any stone this blink from affecting
	// another stone this blink. (eg: if stone=1 is processed first, it'll increase stone=0, then when stone=0 is processed, we only want to
	// process the ones that were there before this blink)
	
	keys := make([]int, 0, len(stones))   // snapshot of stone values (number on the stone)
	counts := make([]int, 0, len(stones)) // snapshot of stone counts (number of stones with that value)

	// Create the snapshot.
	for k := range stones {
		if (stones)[k] > 0 {
			keys = append(keys, k)
			counts = append(counts, stones[k])
		}
	}

	// Iterate over the snapshot.
	for i, stone := range keys {
		count := counts[i]

		// RULE 1
		if stone == 0 {
			(stones)[stone] -= count
			(stones)[1] += count
			continue
		}

		// RULE 2
		digits := countDigits(stone)
		if digits%2 == 0 {
			div := int(math.Pow10(digits / 2))
			leftHalf := stone / div
			rightHalf := stone % div

			(stones)[stone] -= count
			(stones)[leftHalf] += count
			(stones)[rightHalf] += count
			continue
		}

		// RULE 3
		(stones)[stone] -= count
		(stones)[stone * 2024] += count
	}

	// To reduce the number of keys for the next blink, delete any entries with a count of 0.
	for stone, count := range stones {
		if count == 0 {
			delete(stones, stone)
		}
	}

	return stones
}

func blink(stones *[]int) {
	/*
		PART-ONE solution
		Mostly brute force. Following the rules, and updating the array to store the updated stones.
		Iterative, and has a few optimizations but it crashes for part-two after running out of virtual memory.
	*/

	for i, stone := range *stones {
		// RULE 1
		if stone == 0 {
			(*stones)[i] = 1 // Update the stone in place: 0 => 1
			continue
		}

		// RULE 2
		digits := countDigits(stone)
		if digits%2 == 0 {
			// Split the number in two
			div := int(math.Pow10(digits / 2))
			leftHalf := stone / div
			rightHalf := stone % div

			(*stones)[i] = leftHalf                  // Update the stone in place: stone => left_half_of_stone
			(*stones) = append((*stones), rightHalf) // append the right half to the end of the array. yes this will require memory reallocation, but i can't think of a better way.
			continue
		}

		// RULE 3
		// Will apply by default if above 2 fail, so no condition
		(*stones)[i] = stone * 2024 // Update stone in place: stone => stone * 2024
	}
}

func countDigits(number int) int {
	// Could've converted to a string and then just returned length, but this is faster and better for large numbers.
	if number == 0 {
		return 1
	}

	if number < 0 {
		number = -number // Force it to be positive, won't make a difference to digit count.
	}

	return int(math.Log10(float64(number))) + 1
}
