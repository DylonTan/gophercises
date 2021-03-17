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

func parseinput(ch chan string) {
	inputreader := bufio.NewReader(os.Stdin)
	input, err := inputreader.ReadString('\n')

	// Handle no user input error
	if err != nil {
		fmt.Println("Please enter an answer")
	}

	// Send trimmed user input back via channel
	ch <- strings.TrimSpace(input)
}

func checkinput(ch chan string, answer string, limitPtr *int) bool {
	select {
	case input := <- ch:
		// Check if user answer is correct
		if (input == answer) {
			return true
		}
		// Break if timer runs out
	case <- time.After(time.Duration(*limitPtr) * time.Second):
		return false
	}

	return false
}

func main() {
	// Correct answer counter
	correctanswers := 0

	// Parse timer duration flag
	limitPtr := flag.Int("limit", 30, "Limit for timer in seconds")

	// Open problems file
	problemsfile, err := os.Open("problems.csv")

	// Handle unable to open file error
	if err != nil {
		fmt.Println("Problems file could not be opened")
	}

	// Parse problems file
	problemsreader := csv.NewReader(problemsfile)

	for {
		problem, err := problemsreader.Read()

		// Break if no more problems 
		if err != nil {
			break
		}

		// Get questions and answers from problem file
		question, answer := problem[0], problem[1]

		// Trim answer
		answer = strings.TrimSpace(answer)

		// Display question to user
		fmt.Println("What is", question, "?")

		// User input channel
		ch := make(chan string)

		// Parse user input 
		go parseinput(ch)

		// Listen for user input
		if checkinput(ch, answer, limitPtr) {
			correctanswers += 1
		}
	}
	
	// Display correct answers counter
	fmt.Println("Answers correct: ", correctanswers)
}

