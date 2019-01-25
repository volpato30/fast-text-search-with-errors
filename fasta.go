package main

import (
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"log"
	"os"
	"strings"
)

func readFastaFile(f *os.File) [][]string {
	var result [][]string
	t := linear.NewSeq("", nil, alphabet.Protein)
	fastaReader := fasta.NewReader(f, t)
	sc := seqio.NewScanner(fastaReader)
	for sc.Next() {
		s := sc.Seq().(*linear.Seq)
		temp := strings.Split(s.Seq.String(), "")
		result = append(result, temp)
	}
	err := sc.Error()
	if err != nil {
		log.Fatal("failed during reading: %v", err)
	}
	return result
}
