package main

import (
	"context"
	"flag"
	"runtime"
	"time"
	"timepass/task"
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

type line struct {
}

type Quiz struct {
	filename string
}

func New(filename string) *Quiz {

}

func runQuiz() {

}

func main() {

}
