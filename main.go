package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/sushant022/gophercises/quiz"
	"github.com/sushant022/gophercises/task"
)

func runTask() {
	n := flag.Int("n", 100, "Number of tasks")
	duration := flag.Duration("d", 1*time.Second, "Duration to run tasks")
	numConsumers := flag.Int("c", runtime.NumCPU(), "Number consumers")
	flag.Parse()
	ts := task.New(*n, task.Duration(*duration), task.NumConsumers(*numConsumers))
	ts.Run(context.Background())
	ts.Report()
}

func runQuiz() {
	filename := flag.String("f", "quiz.csv", "csv file to load quiz")
	duration := flag.Duration("d", 10*time.Second, "quiz duration")
	flag.Parse()
	q, err := quiz.New(*filename, *duration)
	if err != nil {
		log.Fatal(fmt.Errorf("quiz: Cannot load quiz %w", err))
	}
	q.Take()
	q.Report()
}

func main() {
	runQuiz()
}
