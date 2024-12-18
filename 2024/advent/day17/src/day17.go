package main

// Too many imports, ik.
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	// "math/bits"
	"os"
	"strconv"
	"strings"
)

// That's what the puzzle called the system lol
type StrangeDevice struct {
	registers Registers
	program   []int
}

type Registers struct {
	A int
	B int
	C int
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

	var filePath string
	if !runExampleInput {
		fmt.Println("Running Test input...")
		filePath = "../input/input.txt"
	} else {
		fmt.Println("Running Example input...")
		filePath = "../input/input.example.txt"
	}

	day17(filePath)
}

func day17(filePath string) {
	device := readInput(filePath)
	fmt.Println(device)

	fmt.Println("Part One: ", partOne(device))
	fmt.Println("Part Two: ", partTwo(device))
	// fmt.Println("Part One: ", partOne(device))
}

func partTwo(device StrangeDevice) int {
	// I tried. I really tried to come up with a general solution. But fuck it.

	// Approach: work from the end to the start, and determine the values at each stage.
	// Noticed on the output opcodes that the only value that actually matters when running furst is register A.
	// Register B and C are determined from A each time the opcodes run, so i can basically ignore those.
	// Also, on each stage, register A will only change one way => modulo 8. This means it removes all but the last
	// 3 bits on each stage.
	// So, working backwards from the end of the correct output.
	// Going to try to guess the last 3 bits of A that result
	// in the corresponding output. That's 8 guesses ([0,7]) for each output(determined from the program that is running,
	// iterating backwards).

	// Going to use bfs to search for the value of register A.
	device.registers.A = 0    // remove the value from the input
	stack := make([][]int, 0) // store possible solutions for each output, along with the index for which output it's at

	stack = append(stack, []int{0, 0}) // Starting state.

	for len(stack) > 0 {
		current := stack[0]
		stack = stack[1:]
		testOutput := device.program[len(device.program)-1-current[0]:]

		// i => brute forcing it. going to try to guess the last 3 bits of the register.
		for i := range 8 {
			device.registers.A = (current[1] << 3) | i
			output := partOne(device)
			if arrayEquals(testOutput, output) {
				if arrayEquals(output, device.program) {
					return device.registers.A
				}
				stack = append(stack, []int{current[0] + 1, device.registers.A})
			}
		}
	}
	return 0

}

func arrayEquals(arr1 []int, arr2 []int) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range len(arr1) {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

func partOne(device StrangeDevice) []int {
	output := make([]int, 0)
	for i := 0; i < len(device.program); { // i is the instruction pointer
		opCode := device.program[i]
		operand := device.program[i+1]

		switch opCode {
		case 0: // adv
			result := device.registers.A / int(math.Pow(2, float64(device.toComboOperand(operand))))
			device.registers.A = result
		case 1: // bxl
			result := device.registers.B ^ operand
			device.registers.B = result
		case 2: // bst
			result := device.toComboOperand(operand) % 8
			device.registers.B = result
		case 3: // jnz
			if device.registers.A != 0 {
				i = operand
				continue
			}
		case 4: // bxc
			result := device.registers.B ^ device.registers.C
			device.registers.B = result
		case 5: // out
			result := device.toComboOperand(operand) % 8
			output = append(output, result)
		case 6: //bdv
			result := device.registers.A / int(math.Pow(2, float64(device.toComboOperand(operand))))
			device.registers.B = result
		case 7: // cdv
			result := device.registers.A / int(math.Pow(2, float64(device.toComboOperand(operand))))
			device.registers.C = result
		}
		i += 2
	}
	return output
}

func (device StrangeDevice) toComboOperand(literalOperand int) int {
	switch literalOperand {
	case 0, 1, 2, 3:
		return literalOperand
	case 4:
		return device.registers.A
	case 5:
		return device.registers.B
	case 6:
		return device.registers.C
	case 7:
		return -1
	default:
		return 0
	}
}

func readInput(filePath string) StrangeDevice {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	device := StrangeDevice{
		program: make([]int, 0),
	}

	// Read the registers(until empty line)
	vals := make([]int, 0, 3)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		splitLine := strings.Split(line, ": ")
		registerVal, err := strconv.Atoi(splitLine[1])
		if err != nil {
			log.Fatal("expecting numerical register value: ", err)
		}
		vals = append(vals, registerVal)
	}

	// Set the registers
	device.registers.A = vals[0]
	device.registers.B = vals[1]
	device.registers.C = vals[2]

	scanner.Scan()
	line := scanner.Text()
	splitLine := strings.Split(line, ": ")
	programBuffer := strings.Split(splitLine[1], ",")
	for _, val := range programBuffer {
		numVal, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("expecting numerical opcodes and operands: ", err)
		}
		device.program = append(device.program, numVal)
	}

	return device
}
