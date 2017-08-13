package main

import (
	"log"
	"regexp"
	"strings"
	"time"
)

func main() {

	input := strings.Repeat("a", 50000)
	regex := "aaa"
	var start time.Time
	var elapsed time.Duration

	// Test: There is a match

	// testA
	start = time.Now()
	testA(input, regex)
	elapsed = time.Since(start)
	log.Printf("indexes matching took  %s", elapsed)

	// testB
	start = time.Now()
	testB(input, regex)
	elapsed = time.Since(start)
	log.Printf("indexes range matching took %s", elapsed)

	// Test: There is not a match
	input = strings.Repeat("b", 50000)

	// testA
	start = time.Now()
	testA(input, regex)
	elapsed = time.Since(start)
	log.Printf("indexes not matching took %s", elapsed)

	// testB
	start = time.Now()
	testB(input, regex)
	elapsed = time.Since(start)
	log.Printf("string not matching took %s", elapsed)
}

func testC(input string, regex string) {
	reg := regexp.MustCompile("^" + regex)
	var results [1000000]int

	for i := 0; i < 1000000; i++ {
		indexes := reg.FindStringIndex(input)
		if indexes == nil || indexes[0] != 0 {
			results[i] = 0
			continue
		}
		results[i] = indexes[1] // length of match
	}
}

func testA(input string, regex string) {
	reg := regexp.MustCompile("^" + regex)
	var results [1000000]int

	for i := 0; i < 1000000; i++ {
		str := reg.FindString(input)
		index := strings.Index(input, str)
		if index != 0 {
			results[i] = 0
			continue
		}
		results[i] = len(str) // length of match
	}
}

func testB(input string, regex string) {
	reg := regexp.MustCompile("^" + regex)
	var results [1000000]int

	for i := 0; i < 1000000; i++ {
		str := reg.FindString(input)
		strlen := len(str)
		if input[0:strlen] != str {
			results[i] = 0
			continue
		}
		results[i] = strlen // length of match
	}
}
