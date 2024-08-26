package main

import (
	"bufio"
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

type SolutionSets map[int]*SolutionSet

func main() {
	fileName, duration := parseFlags()
	solutionSets, err := prepareSolutionSets(fileName)
	if err != nil {
		fmt.Printf("Error preparing solution sets: %v\n", err)
		os.Exit(1)
	}
	startTest(solutionSets, duration)
	generateResult(solutionSets)
}

func parseFlags() (string, int) {
	fileName := flag.String("filename", "quiz.csv", ".csv Answer Key Filepath")
	duration := flag.Int("duration", 10, "Duration for test in seconds")
	flag.Parse()
	if *fileName == "" {
		fmt.Println("Please provide Question Answer Filename")
		os.Exit(1)
	}
	return *fileName, *duration
}

func takeTest(solutionSets SolutionSets, done chan<- bool) {
	reader := bufio.NewReader(os.Stdin)
	for _, v := range solutionSets {
		fmt.Printf("%v: ", v.Question)
		userResponse, _ := reader.ReadString('\n')
		v.UserResponse = strings.TrimSpace(userResponse)
		v.IsCorrect = strings.EqualFold(v.UserResponse, v.Answer)
	}
	done <- true
}

func prepareSolutionSets(fileName string) (SolutionSets, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open quiz file: %w", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read quiz file: %w", err)
	}

	solutionSets := make(SolutionSets)
	for i, record := range records {
		if len(record) != 2 {
			return nil, fmt.Errorf("invalid record format at line %d", i+1)
		}
		solutionSets[i] = &SolutionSet{
			Question: strings.TrimSpace(record[0]),
			Answer:   strings.TrimSpace(record[1]),
		}
	}
	return solutionSets, nil
}

func startTest(solutionSets SolutionSets, duration int) {
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	done := make(chan bool)
	go takeTest(solutionSets, done)

	select {
	case <-timer.C:
		fmt.Println("\nTime is up")
	case <-done:
		timer.Stop()
	}
}

func generateResult(solutionSets SolutionSets) {
	correctAttempts := 0
	for _, v := range solutionSets {
		if v.IsCorrect {
			correctAttempts++
		}
	}
	percentage := float64(correctAttempts) / float64(len(solutionSets)) * 100
	fmt.Printf("You have scored %.2f%%\n", percentage)
}
