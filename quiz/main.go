package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type SolutionSet struct {
	Question     string
	Answer       string
	UserResponse string
	IsCorrect    bool
}

type UserSet struct {
	*SolutionSet
	UserResponse string
	IsCorrect    bool
}

var (
	fileName string
	duration int
)

type SolutionSets map[int]*SolutionSet

func dirationParser() {

}

func init() {
	flag.StringVar(&fileName, "filename", "quiz.csv", ".csv Answer Key Filepath")
	flag.IntVar(&duration, "duration", 10, "Duration for test in seconds")
	flag.Parse()
	if fileName == "" {
		panic("Please provide Question Answer Filename")
	}
}

func takeTest(solutionSets SolutionSets, done chan<- bool) {
	for _, v := range solutionSets {
		var userResponse string
		fmt.Printf("%v:", v.Question)
		fmt.Scanln(&userResponse)

		userResponse = strings.TrimSpace(userResponse)
		v.UserResponse = userResponse
		if userResponse == v.Answer {
			v.IsCorrect = true
		}
	}
	done <- true
}

func prepareSolutionSets(fileName string) SolutionSets {
	f, err := os.Open(fileName)
	if err != nil {
		panic("cannot read quiz file")
	}
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic("cannot read quiz file")
	}

	solutionSets := map[int]*SolutionSet{}

	for i, record := range records {
		solutionSet := SolutionSet{
			Question: strings.TrimSpace(record[0]),
			Answer:   strings.TrimSpace(record[1]),
		}
		solutionSets[i] = &solutionSet
	}
	return solutionSets
}

func startTest(solutionSets SolutionSets) {
	timer := time.NewTimer(time.Duration(duration) * time.Second).C
	done := make(chan bool)
	go takeTest(solutionSets, done)

loop:
	for {
		select {
		case <-timer:
			fmt.Println("\nTime is up")
			break loop
		case <-done:
			break loop
		}
	}
}

func generateResult(solutionSets SolutionSets) {
	var (
		correctAttempts int
		percentage      float64
	)
	for _, v := range solutionSets {
		if v.UserResponse == v.Answer {
			correctAttempts++
		}
	}
	percentage = float64(correctAttempts) / float64(len(solutionSets)) * 100
	fmt.Printf("You have scored %.2f%%\n", percentage)
}

func main() {
	// Read csv and prepare Question Answer set
	solutionSets := prepareSolutionSets(fileName)

	// Take Test
	startTest(solutionSets)

	// Evaluate Results
	generateResult(solutionSets)
}
