package main

// Too many imports, ik.
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sort"
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

	day12(filePath, 1, 2)
}

func day12(filePath string, parts ...int) {
	input := readInput(filePath)
	for _, part := range parts {
		switch part {
		case 1:
			part1(input)
		case 2:
			part2(input)
		}
	}
}

func part1(input [][]string) {
	// Adjacency matrix basically. Connections will be one sided tho, so check both directions for any given connection.
	adj := make(map[string]map[string]bool)
	for _, connection := range input {
		// Form the adjacency matrix
		if adj[connection[1]][connection[0]] { // repeat connection (just in case)
			continue
		}

		if _, exists := adj[connection[0]]; !exists {
			adj[connection[0]] = make(map[string]bool)
		}
		if _, exists := adj[connection[1]]; !exists {
			adj[connection[1]] = make(map[string]bool)
		}

		adj[connection[0]][connection[1]] = true
		adj[connection[1]][connection[0]] = true
	}

	lanParties := [][3]string{}

	visited := make(map[string]bool) // for pc 1
	for pc1 := range adj {
		if visited[pc1] {
			continue
		}
		visited[pc1] = true

		// for pc2 and pc3
		_visited := make(map[string]bool) // order doesn't matter so.
		for pc2 := range adj[pc1] {
			if _visited[pc2] || visited[pc2] {
				continue
			}
			_visited[pc2] = true
			for pc3 := range adj[pc2] {
				if _visited[pc3] || visited[pc3] {
					continue
				}

				if adj[pc1][pc3] {
					lanParties = append(lanParties, [3]string{pc1, pc2, pc3})
				}
			}
		}
	}
	fmt.Println(lanParties)

	count := 0
	for _, party := range lanParties {
		for _, pc := range party {
			if pc[0] == 't' {
				count++
				break
			}
		}
	}

	fmt.Println("Part One: ", count)

}

func part2(input [][]string) {
	// Graph theory. yay.
	adj := make(map[string]map[string]bool)
	for _, connection := range input {
		// Form the adjacency matrix
		if adj[connection[1]][connection[0]] { // repeat connection (just in case)
			continue
		}

		if _, exists := adj[connection[0]]; !exists {
			adj[connection[0]] = make(map[string]bool)
		}
		if _, exists := adj[connection[1]]; !exists {
			adj[connection[1]] = make(map[string]bool)
		}

		adj[connection[0]][connection[1]] = true
		adj[connection[1]][connection[0]] = true
	}

	// Find the maximal cliques
	cliques := make([][]string, 0)
	P := make([]string, 0)
	for v := range adj {
		P = append(P, v)
	}

	bron_kerbosch(adj, []string{}, P, []string{}, &cliques)

	fmt.Println("Part Two")
	
	maxClique := []string{}
	maxCliqueLength := -1

	for _, clique := range cliques {
		if len(clique) > maxCliqueLength {
			maxCliqueLength = len(clique)
			maxClique = clique
		}
	}

	sort.Strings(maxClique)
	fmt.Println("Length: ", maxCliqueLength)
	fmt.Println("Clique: ", maxClique)
}

func bron_kerbosch(adj map[string]map[string]bool, R []string, P []string, X []string, cliques *[][]string) {
	// Basic Bron-Kerbosh algorithm for finding the maximal cliques. This will work through them all and find
	// all cliques in the given adjacency matrix. 

	// Base case - if both P and X are empty and R is not empty, we found a maximal clique
	if len(P) == 0 && len(X) == 0 && len(R) > 0 {
		maximalClique := make([]string, len(R))
		copy(maximalClique, R)
		*cliques = append(*cliques, maximalClique)
		return
	}

	// If P is empty, we can't expand further
	if len(P) == 0 {
		return
	}

	// For each vertex in P
	for i := 0; i < len(P); i++ {
		v := P[i]

		// Create new R by adding v
		newR := append([]string{}, R...)
		newR = append(newR, v)

		// Create new P and X by intersecting with neighbors of v
		newP := make([]string, 0)
		for _, p := range P[i+1:] {
			if adj[v][p] {
				newP = append(newP, p)
			}
		}

		newX := make([]string, 0)
		for _, x := range X {
			if adj[v][x] {
				newX = append(newX, x)
			}
		}

		bron_kerbosch(adj, newR, newP, newX, cliques)

		// Move v from P to X
		X = append(X, v)
	}
}

func readInput(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	buffer := make([][]string, 0, 1024)
	// Simply reading the input. No transformations for this stage
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "-")
		buffer = append(buffer, line)
	}
	return buffer
}
