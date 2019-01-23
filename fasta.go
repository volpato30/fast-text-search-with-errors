package main

import (
	"bufio"
	"io"
	"strings"
)

type FastaRecord struct {
	header   string
	sequence string
}

func NewFastxReader(f io.Reader) *FastxReader {
	return &FastxReader{
		r: bufio.NewReader(f),
	}
}

type FastxReader struct {
	r *bufio.Reader
}

func (r *FastxReader) next_seq() (record FastaRecord, err error) {
	var str string
	if str, err = r.r.ReadString('>'); err == nil {
		if str, err = r.r.ReadString('>'); err == nil {
			split_result := strings.SplitN(str, "\n", 2)
			record.header = split_result[0]
			//remove newlines and trailing >
			record.sequence = chomp(strings.Replace(split_result[1], "\n", "", -1), ">")
		}
	}
	return record, err
}

//remove last char in a string if that char is the delim
func chomp(s string, delim string) string {
	if s[len(s)-1] == delim[0] {
		return s[0 : len(s)-1]
	}
	return s
}
