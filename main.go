package main

import (
	"fmt"
//	"sync"
	"log"
)

//var fileMutex sync.Mutex

func test_hamming1_search() {
	pattern := []string{"T", "C", "G", "T"}
	text := []string{"T",
		"T",
		"T",
		"A",
		"C",
		"G",
		"T",
		"A",
		"A",
		"A",
		"C",
		"T",
		"A",
		"A",
		"A",
		"C",
		"T",
		"G",
		"T",
		"A",
		"A"}
	matchLoc, err := HammingDist1Search(pattern, text, StringIdentityMap)
	if err != nil {
		log.Fatal(err)
	}

	if matchLoc.exactMatch {
		panic("test not passed")
	}
	if len(matchLoc.startIndex) != 1 {
		panic("test not passed")
	}
	if matchLoc.startIndex[0] != 3 {
		panic("test not passed")
	}
	pattern = []string{"A", "A", "C", "T"}
	matchLoc, err = HammingDist1Search(pattern, text, StringIdentityMap)
	if err != nil {
		log.Fatal(err)
	}
	if !matchLoc.exactMatch {
		panic("test not passed")
	}

}

func main() {
	test_hamming1_search()
	fmt.Println("test passed")
}

