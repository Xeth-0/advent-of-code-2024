package main

import (
	"log"
	"strconv"
	"bufio"
	"os"
	"regexp"
)

func main () {
	day3()
}

func day3() {
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal("error reading file")
	}
	defer file.Close()

	// part one
	// sum := 0
	scanner := bufio.NewScanner(file)
	buffer := ""
	for scanner.Scan() {
		line := scanner.Text()
		buffer += line + "\n"
	}

	// find the matches(mul(\d,\d)) for the input
	log.Println(buffer)
	// buffer = "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"
	re := regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d+),(\d+)\)`)
	matches := re.FindAllStringSubmatch(buffer, -1)
	log.Println("matches: ", matches)

	// sum up the matches
	sum := 0
	enableMul := true
	for _, match := range matches {
		if match[0] == "don't()" {
			enableMul = false
			continue
		} else if (match[0]) == "do()" {
			enableMul = true
			continue
		}

		if enableMul {
			num1, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal("error converting match to string: ", match, err)
			}
			num2, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal("error converting match to string: ", match, err)
			}
	
			sum += num1 * num2
		}
	}

	log.Println("Sum of matches is: ", sum)
}
