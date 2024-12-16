package main

import (
	// "fmt"
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	solution()
}

func solution() {
	// read the file
	log.Println("reading file...")
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	defer file.Close()

	// stores for the two lists of numbers
	list1 := make([]int, 0, 512)
	list2 := make([]int, 0, 512)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Split(line, "   ")
		num1, err := strconv.Atoi(numbers[0])
		if err != nil {
			log.Fatal("error parsing int: ", err)
		}
		num2, err := strconv.Atoi(numbers[1])
		if err != nil {
			log.Fatal("error parsing int: ", err)
		}

		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}
	sort.Ints(list1)
	sort.Ints(list2)

	// Calculate the distance
	distance := 0
	for i := range len(list1) {
		d := list1[i] - list2[i]
		if d < 0 {
			d *= -1
		}
		distance += d
	}

	log.Println("PART-ONE SOLUTION: ", distance)

	// Part-TWO
	i := 0
	j := 0

	similarityScore := 0
	for i < len(list1) && j < len(list2){
		if list1[i] < list2[j] {
			i++
		} else if list1[i] == list2[j] {
			similarityScore += list1[i]
			j++
		} else {
			j++
		}
	}

	log.Println("PART-TWO SOLUTION: ", similarityScore)
}
