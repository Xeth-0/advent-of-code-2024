package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Parse args
	runExampleInputFlag := flag.Bool("e", true, "run the example input or the test input")
	flag.Parse()

	runExampleInput := *runExampleInputFlag

	// Set up logger
	if runExampleInput {
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

	day9(filePath)
}

func day9(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([]int, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		for _, v := range line {
			val, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("error parsing int: %v", err)
			}

			buffer = append(buffer, val)
		}
	}

	log.Println("Original Input: ", buffer)
	disk := parseDiskMap(buffer)
	log.Println("Compacted Disk: ", disk)

	// Part One
	// Compact the disk
	compactedDisk := compactDisk(copyDisk(disk))
	checksum := chksum(compactedDisk)
	fmt.Println("CHECKSUM: ", checksum)

	compactedDisk = compactDiskByBlock(copyDisk(disk))
	fmt.Println("Compacted Disk: ", compactedDisk)
	checksum = chksum(compactedDisk)
	fmt.Println("CHECKSUM: ", checksum)
}

func chksum(disk []int) int {
	ret := 0

	for i, val := range disk {
		if val == -1 {
			continue
		}
		ret += (i * val)
	}

	return ret
}

func compactDisk(disk []int) []int {
	i := 0
	j := len(disk) - 1

	for i < len(disk) && i < j && j > 0 {
		for disk[i] != -1 {
			i++
		}
		for disk[j] == -1 {
			j--
		}

		disk[i] = disk[j]
		disk[j] = -1

		i++
		j--
	}
	return disk
}

func compactDiskByBlock(disk []int) []int {
	// Map out the free space
	freeSpaces := make(map[int]int)
	freeSpaceIndexes := make([]int, 0)
	for i := 0; i < len(disk); {
		if disk[i] == -1 {
			freeStart := i
			for disk[i] == -1 {
				i++
			}
			freeSpaces[freeStart] = i - freeStart // length of free space
			freeSpaceIndexes = append(freeSpaceIndexes, freeStart)
		} else {
			i++
		}
	}
	// fmt.Println(freeSpaces)

	// Compaction
	for i := len(disk) - 1; i >= 0; i-- {
		// Find the end of a file.
		for disk[i] == -1 {
			i--
		}

		// Iterate to the start of the file block. Should already be at the end
		blockEnd := i
		for i > 0 && disk[i] == disk[blockEnd] {
			i--
		}
		if i <= 0 {
			break
		}

		// Get the block of files.
		i++
		block := disk[i : blockEnd+1]

		// Assign it to a contiguous free block.
		for k, freeIndex := range freeSpaceIndexes {
			freeLength := freeSpaces[freeIndex]
			if i > freeIndex && len(block) > 0 && freeLength >= len(block) {
				for j, val := range block {
					disk[freeIndex+j] = val
					block[j] = -1
				}
				// Update the free space identifiers.
				freeSpaces[freeIndex] = 0
				if freeLength >= len(block) { // free space left.
					freeSpaces[freeIndex+len(block)] = freeLength - len(block)
					freeSpaceIndexes[k] = freeIndex + len(block)
				}
				break
			}
		}
	}
	return disk
}

func copyDisk(disk []int) []int {
	copy := make([]int, 0, len(disk))
	copy = append(copy, disk...)
	return copy
}

func parseDiskMap(disk []int) []int {
	updatedDisk := make([]int, 0, len(disk))
	fileID := 0

	for i, val := range disk {
		if i%2 == 0 { // File
			for range val {
				updatedDisk = append(updatedDisk, fileID)
			}
			fileID++
		} else { // Space
			for range val {
				updatedDisk = append(updatedDisk, -1)
			}
		}
	}

	return updatedDisk
}
