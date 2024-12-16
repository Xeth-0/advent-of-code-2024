package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	day5()
}

func day5() {
	// Parse the input
	// Setup
	filePath := "../input/input.txt"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	scanner := bufio.NewScanner(file)

	// Read the first section: the update rules
	log.Println("Reading rules...")
	rules := make(map[int][]int, 0) // map to store the rules for the updates.
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // Empty line denotes the end of the first section.
			break
		}
		newRule := strings.Split(line, "|")

		num1, err := strconv.Atoi(newRule[0])
		if err != nil {
			log.Fatal("Non-numerical rule found: ", err)
		}

		num2, err := strconv.Atoi(newRule[1])
		if err != nil {
			log.Fatal("Non-numerical rule found: ", err)
		}

		// let's add the rule to the rules map
		rule, exists := rules[num2]
		if !exists {
			rules[num2] = make([]int, 0, 1)
			rule = rules[num2]
		}

		rule = append(rule, num1)
		rules[num2] = rule
	}
	log.Println("Rules: ", rules)

	// Second section: Validating updates
	log.Println("Reading updates...")

	count := 0  // answer for part 1
	count2 := 0 // answer for part 2

	for scanner.Scan() {
		rawUpdate := strings.Split(scanner.Text(), ",")
		newUpdate := make([]int, 0, len(rawUpdate))

		// convert them into ints (could be pointless but whatever)
		for _, val := range rawUpdate {
			numVal, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal("Got a non-numerical update value: ", val, err)
			}
			newUpdate = append(newUpdate, numVal)
		}

		isValid := validateUpdate(newUpdate, rules)
		if isValid {
			mid := 0
			l := len(newUpdate)
			if l > 1 {
				if l%2 != 0 {
					l++
				}
				mid = (l / 2) - 1
			}
			count += newUpdate[mid]
		}

		// Part-Two
		if !isValid { // for the invalid updates only
			newUpdate = correctUpdate(newUpdate, rules)
			mid := 0
			l := len(newUpdate)
			if l > 1 {
				if l%2 != 0 {
					l++
				}
				mid = (l / 2) - 1
			}
			count2 += newUpdate[mid]

		}
	}
	log.Println("DAY5-PART1: ", count)
	log.Println("DAY5-PART2: ", count2)
}

// Validates given update (list of ints) and returns (isValid, the middle update).
func validateUpdate(update []int, rules map[int][]int) (isValid bool) {
	// Let's map the update page values to their indexes for easier validation
	valueIndexes := make(map[int]int, 0)
	for i, val := range update {
		valueIndexes[val] = i // does not handle the same page occurring more than once
	}

	for i, val := range update {
		// Let's validate the pages using the rules
		rule, exists := rules[val]
		if !exists {
			continue
		}
		for _, precursor := range rule {
			idx, exists := valueIndexes[precursor]
			if exists && idx > i {
				return false
			}
		}
	}
	return true
}

func correctUpdate(update []int, rules map[int][]int) (correctedUpdate []int) {
	// first let's check if the update is valid
	if validateUpdate(update, rules) {
		return update
	}

	// Let's map the update page values to their indexes for easier validation
	valueIndexes := make(map[int]int, 0)
	for i, val := range update {
		valueIndexes[val] = i // does not handle the same page occurring more than once
	}

	for i, val := range update {
		// Let's validate the pages using the rules
		rule, exists := rules[val]
		if !exists {
			continue
		}
		for _, precursor := range rule {
			idx, exists := valueIndexes[precursor]
			if exists && idx > i {
				// Found what was making it invalid
				correctedUpdate = make([]int, 0, len(update))
				correctedUpdate = append(correctedUpdate, update[:i]...)
				correctedUpdate = append(correctedUpdate, update[idx])
				correctedUpdate = append(correctedUpdate, update[i:idx]...)
				if idx+1 < len(update) {
					correctedUpdate = append(correctedUpdate, update[idx+1:]...)
				}

				return correctUpdate(correctedUpdate, rules)
			}
		}
	}
	return update
}
