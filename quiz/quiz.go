package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type line struct {
	question    string
	answer      string
	response    string
	isCorrect   bool
	isAttempted bool
}

type Quiz struct {
	lines     []*line
	duration  time.Duration
	Correct   int
	Attempted int
}

// takes filename as input and prepares questions and answers for quiz
func New(filename string, duration time.Duration) (*Quiz, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	q := Quiz{duration: duration}
	for _, r := range rows {
		q.lines = append(q.lines, &line{question: strings.TrimSpace(r[0]), answer: strings.TrimSpace(r[1])})
	}
	return &q, nil
}

func ask() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var answer string
		fmt.Scanln(&answer)
		ch <- strings.TrimSpace(answer)
	}()
	return ch
}

func (q *Quiz) Take() {
	timeout := time.After(q.duration)
	for i, l := range q.lines {
		fmt.Printf("%d. %s:\n", i+1, l.question)
		answerCh := ask()
		select {
		case <-timeout:
			return
		case v := <-answerCh:
			l.response = v
			if l.response == l.answer {
				l.isCorrect = true
				q.Correct++
			}
			l.isAttempted = true
			q.Attempted++
		}
	}
}

func boolToStr(b bool) string {
	if !b {
		return "N"
	}
	return "Y"
}

func (q *Quiz) Report() {
	fmt.Printf("\n\t\tREPORT\t\t\n\n")
	fmt.Println("Question | Answer | Response | Is correct?")
	for _, l := range q.lines {
		if l.isAttempted {
			fmt.Printf("%8s | %6s | %8s | %11s\n", l.question, l.answer, l.response, boolToStr(l.isCorrect))
		}
	}
	fmt.Printf("\n\t\tSUMMARY\t\t\n\n")
	fmt.Println("Correct | Attempted | Total | Percentage")
	fmt.Printf("%7d | %9d | %5d | %10.2f%%\n", q.Correct, q.Attempted, len(q.lines), float64(q.Correct)/float64(len(q.lines))*100)
}
