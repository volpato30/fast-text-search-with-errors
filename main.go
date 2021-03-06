package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Configuration struct {
	denovoPeptideFilename string
	fastaFileName         string
	outputFileName        string
	batchSize             int
}

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func batchProcess(denovoPeptides []string, databasePeptides [][]string, writer *bufio.Writer, mutex *sync.Mutex,
	wg *sync.WaitGroup) {
	defer wg.Done()
	var batchOutputStrings []string
	for _, peptide := range denovoPeptides {
		peptide = strings.Trim(peptide, "\n")
		aa_seq := strings.Split(peptide, "")
		peptideLength := uint(len(aa_seq))
		var wildTypePeptides []string
		for ii, _ := range databasePeptides {
			referenceSeq := databasePeptides[ii]
			matchLoc, err := HammingDist1Search(aa_seq, referenceSeq, IToLMap)
			check(err)
			if matchLoc.exactMatch {
				continue
			} else if len(matchLoc.startIndex) == 0 {
				continue
			} else {
				for _, m := range matchLoc.startIndex {
					wt := strings.Join(referenceSeq[m:m+peptideLength], "")
					wildTypePeptides = append(wildTypePeptides, wt)
				}

			}
		}
		if len(wildTypePeptides) > 0 {
			outputString := strings.Join(wildTypePeptides, ";")
			outputString = peptide + "\t" + outputString + "\n"
			batchOutputStrings = append(batchOutputStrings, outputString)
		}
	}
	mutex.Lock()
	defer mutex.Unlock()
	for ii, _ := range batchOutputStrings {
		_, err := writer.WriteString(batchOutputStrings[ii])
		check(err)
		err = writer.Flush()
		check(err)
	}

}

func main() {
	test_hamming1_search()
	fmt.Println("test passed")
	config := Configuration{
		denovoPeptideFilename: "./peptide1.txt",
		fastaFileName:         "wild_type_match/fasta_files/HUMAN.fasta",
		outputFileName:        "wild_type_match/MSV000097215/wt_aligned.peptide.txt",
		batchSize:             200,
	}
	start := time.Now()
	f, err := os.Open(config.fastaFileName)
	defer f.Close()
	check(err)
	databasePeptides := readFastaFile(f)
	fmt.Printf("read %d proteins\n", len(databasePeptides))
	var wg sync.WaitGroup
	var denovoPeptides []string
	var mutex sync.Mutex

	denovoFile, err := os.Open(config.denovoPeptideFilename)
	check(err)
	defer denovoFile.Close()
	outputFile, err := os.Create(config.outputFileName)
	check(err)
	writer := bufio.NewWriter(outputFile)
	defer outputFile.Close()

	scanner := bufio.NewScanner(denovoFile)
	for scanner.Scan() {
		temp := scanner.Text()
		denovoPeptides = append(denovoPeptides, temp)
		if len(denovoPeptides) >= config.batchSize {
			wg.Add(1)
			go batchProcess(denovoPeptides, databasePeptides, writer, &mutex, &wg)
			fmt.Println("start a goroutine")
			denovoPeptides = []string{}
		}
	}

	wg.Add(1)
	go batchProcess(denovoPeptides, databasePeptides, writer, &mutex, &wg)
	fmt.Println("start a goroutine")
	denovoPeptides = []string{}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
