package main

import (
	"fmt"
)

type MatchLocations struct {
	exactMatch bool
	startIndex []uint
}

func StringIdentityMap(char string) string {
	return char
}

func IToLMap(char string) string {
	if char == "I" {
		return "L"
	} else {
		return char
	}
}

type CharConvert func(string) string

type argError struct {
	message string
}

func (e *argError) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func HammingDist1Search(pattern []string, text []string, convertFunc CharConvert) (MatchLocations, error) {
	m := uint(len(pattern))
	var matchLoc MatchLocations
	if m > 30 {
		return matchLoc, &argError{"do not support pattern longer than 30"}
	}
	R0, R1, S := 0, 0, 0
	var shR0 int
	var mappedChar string
	mask := 1 << (m - 1)

	sMap := make(map[string]int)

	for index, char := range pattern {
		mappedChar = convertFunc(char)
		val, ok := sMap[mappedChar]
		if ok {
			sMap[mappedChar] = val | (1 << uint(index))
		} else {
			sMap[mappedChar] = 1 << uint(index)
		}
	}

	for index, char := range text {
		mappedChar = convertFunc(char)
		val, ok := sMap[mappedChar]
		if ok {
			S = val
		} else {
			S = 0
		}
		shR0 = (R0 << 1) | 1
		R0 = shR0 & S
		R1 = (((R1 << 1) | 1) & S) | shR0
		if R0&mask != 0 {
			matchLoc.exactMatch = true
			return matchLoc, nil
		} else if R1&mask != 0 {
			matchLoc.startIndex = append(matchLoc.startIndex, uint(index-int(m)+1))
		}
	}
	return matchLoc, nil
}
