// This system will take in a number of tasks and proccess them in batches of 5 at a time.
package task

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"
)

type TaskProcessor struct {
	tasks        []*Task
	numTasks     int
	numConsumers int
	duration     time.Duration
	isverbose    bool
	completed    int
}

type Option func(t *TaskProcessor)

func NumConsumers(n int) Option {
	return func(ts *TaskProcessor) {
		ts.numConsumers = n
	}
}

func isverbose(v bool) Option {
	return func(t *TaskProcessor) {
		t.isverbose = v
	}
}

func Duration(d time.Duration) Option {
	return func(ts *TaskProcessor) {
		ts.duration = d
	}
}

func New(numTasks int, options ...Option) *TaskProcessor {
	t := TaskProcessor{tasks: generateTasks(numTasks), numTasks: numTasks, numConsumers: runtime.NumCPU()}
	for _, opt := range options {
		opt(&t)
	}
	return &t
}

func generateTasks(n int) []*Task {
	tasks := make([]*Task, n)

	for i := 0; i < n; i++ {
		tasks[i] = &Task{
			Id:       i,
			Duration: time.Duration(rand.IntN(1000)) * time.Millisecond,
		}
	}
	return tasks
}

func (ts *TaskProcessor) Run(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, ts.duration)
	defer cancel()
	producerStream := produce(ctx, ts.tasks)
	for range mux(ctx, ts.numConsumers, producerStream) {
		ts.completed++
	}
}

func (ts TaskProcessor) Report() {
	fmt.Printf("completed: %d\ntotal: %d\npercentage:%.2f%%\n", ts.completed, len(ts.tasks), float64(ts.completed)/float64(len(ts.tasks))*100)
}

type Task struct {
	Id          int
	Duration    time.Duration
	Iscompleted bool
}

func produce(ctx context.Context, tasks []*Task) <-chan *Task {
	ch := make(chan *Task)
	go func() {
		defer close(ch)
		for _, t := range tasks {
			select {
			case <-ctx.Done():
				return
			case ch <- t:
			}
		}
	}()
	return ch
}

func consumer(ctx context.Context, inStream <-chan *Task) <-chan *Task {
	ch := make(chan *Task)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-inStream:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(t.Duration):
					fmt.Printf("Id: %d, Status: processing, Duration: %d\n", t.Id, t.Duration)
					t.Iscompleted = true
					ch <- t
				}
			}
		}
	}()
	return ch
}

func mux(ctx context.Context, numConsumers int, producerStream <-chan *Task) <-chan *Task {
	fanInStream := make(chan *Task)
	consumerStreams := make([]<-chan *Task, numConsumers)

	for i := 0; i < numConsumers; i++ {
		consumerStreams[i] = consumer(ctx, producerStream)
	}

	var wg sync.WaitGroup
	wg.Add(numConsumers)

	for i := 0; i < numConsumers; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-consumerStreams[i]:
					if !ok {
						return
					}
					fanInStream <- v
				}
			}
		}(i, &wg)
	}

	go func(wg *sync.WaitGroup) {
		defer close(fanInStream)
		wg.Wait()
	}(&wg)

	return fanInStream
}
