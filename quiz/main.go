package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "A cvs file containing the math problems.")
	timeLimit := flag.Int("limit", 30, "Define time limit for each question.")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: s%", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: s%", *csvFileName))
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	result := 0

	fmt.Printf("You have %s to answer each question! \n", time.Duration(*timeLimit)*time.Second)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s ", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			var userInput string
			fmt.Scanf("%s\n", &userInput)
			answerCh <- userInput
		}()

		select {
		case <-timer.C:
			fmt.Printf("You have made %d of %d \n", result, len(problems))
			return
		case userInput := <-answerCh:
			if userInput == p.answer {
				result++
			}
		}

	}
	fmt.Printf("You have made %d of %d \n", result, len(problems))

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	println(msg)
	os.Exit(1)
}
