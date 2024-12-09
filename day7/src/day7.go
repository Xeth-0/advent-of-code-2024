package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// day7("../input/input.example.txt")
	day7("../input/input.txt")
}

func day7(filePath string) {
	// Read the input
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("cannot read input file: ", err)
	}
	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ": ")
		testValue, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal("expecting testValue to be numerical: ", err)
		}
		values := strings.Split(line[1], " ")
		nums := make([]int, len(values))

		for i, val := range values {
			v, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal("expecting nums to be numerical: ", err)
			}
			nums[i] = v
		}

		// fmt.Println(nums, testValue)
		if checkTestValue(nums[0], nums[1:], testValue) {
			count += testValue
		}
	}
	fmt.Println(count)

}

func checkTestValue(previousVal int, values []int, testValue int) bool {
	if len(values) == 0 {
		return previousVal == testValue
	} else if len(values) == 1 {
		return previousVal+values[0] == testValue || previousVal*values[0] == testValue || concatenateInts(previousVal, values[0]) == testValue
	}

	return checkTestValue(values[0]+previousVal, values[1:], testValue) ||
		checkTestValue(values[0]*previousVal, values[1:], testValue) ||
		checkTestValue(concatenateInts(previousVal, values[0]), values[1:], testValue)
}

func concatenateInts(num1 int, num2 int) int {
	val, _ := strconv.Atoi(fmt.Sprintf("%d%d", num1, num2))
	return val
}

// Recursive solution map:
/*
	i need the recursor to try all possible operators.
	at each "STAGE" of the recursive loop, we have to choose which operator should be used proceed,
	[12 | 3 4 5]
	answer is 3 * 4 + 5
	it's evaluated Left to Right, so i dont need to be worried about precedence. Whatever the previous value was will be fed into the next operator.



*/
