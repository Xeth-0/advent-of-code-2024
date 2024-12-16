package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	day2()
}

func day2(){
	log.Println("reading file ...")
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	defer file.Close()

	// Prase the file
	data := make([][]int, 0, 1024)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dataPoints := strings.Split(line, " ")
		level := make([]int, 0, len(dataPoints))
		for _, datapoint := range dataPoints {
			i, err := strconv.Atoi(datapoint)
			if err != nil {
				log.Fatal("error reading file: unexpected format. Expected ints separated by ' ': ", err)
			}
			level = append(level, i)
		}

		data = append(data, level)
	}

	// part one
	numSafe := 0
	for _, level := range data {
		if checkSafe(level, 1) {
			numSafe++
			log.Println("SAFE")
		}
	}
	log.Println("NUMSAFE: ", numSafe)
}

func checkSafe(level []int, tolerance int) bool {
	if len(level) == 1 {
		return true
	}
	if tolerance < 0 {
		return false
	}

	decreasing := level[0] > level[1]
	log.Println(level)
	// log.Println(decreasing)
	for i := range len(level) - 1 {
		if decreasing {
			if level[i] < level[i+1] || level[i]-level[i+1] > 3 || level[i] == level[i+1] {
				log.Println(level[i], level[i+1])

				if i >= 1 {
					excludePrev := append([]int(nil), level[:i-1]...)
					excludePrev = append(excludePrev, level[i:]...)

					if checkSafe(excludePrev, tolerance - 1) {
						return true
					}
				}

				excludeFirst := append([]int(nil), level[:i]...) 
				excludeFirst = append(excludeFirst, level[i+1:]...)
				
				excludeSecond := append([]int(nil), level[:i+1]...)
				excludeSecond = append(excludeSecond, level[i+2:]...)

				return checkSafe(excludeFirst, tolerance - 1) || checkSafe(excludeSecond, tolerance - 1)
			}
		} else {
			if level[i] > level[i+1] || level[i+1]-level[i] > 3 || level[i] == level[i+1] {
				log.Println(level[i], level[i+1])

				if i >= 1 {
					excludePrev := append([]int(nil), level[:i-1]...)
					excludePrev = append(excludePrev, level[i:]...)

					if checkSafe(excludePrev, tolerance - 1) {
						return true
					}
				}


				excludeFirst := append([]int(nil), level[:i]...) 
				excludeFirst = append(excludeFirst, level[i+1:]...)
				
				excludeSecond := append([]int(nil), level[:i+1]...)
				excludeSecond = append(excludeSecond, level[i+2:]...)

				return checkSafe(excludeFirst, tolerance - 1) || checkSafe(excludeSecond, tolerance - 1)
			}
		}
	}
	return true
}
