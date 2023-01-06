package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func main() {
	fileName := flag.String(
		"csv",
		"problems.csv",
		"A CSV file for the quiz content in question,answer format.")
	timeLimit := flag.Int("limit", 30, "The time limit for the quiz.")
	shuffle := flag.Bool(
		"randomize",
		false,
		"Should the questions be presented in random order.")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("The provided file [%s] couldn't be opened.\n", *fileName)
		fmt.Println(err)
		os.Exit(1)
	}

	problems := parseLines(file)
	if *shuffle {
		problems = shuffleProblems(problems)
	}

	fmt.Print("Press enter to begin...")
	fmt.Scanf("\n")
	fmt.Println()

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	done := make(chan struct{})
	correct := 0
	go func() {
		for i, prob := range problems {
			fmt.Printf("Question #%d: %s, Your answer=", i+1, prob.question)
			var answer string
			fmt.Scanf("%s\n", &answer)
			if answer == prob.answer {
				correct++
			}
		}
		timer.Stop()
		close(done)
	}()

	select {
	case <-timer.C:
	case <-done:
	}

	fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
}

func parseLines(reader io.Reader) []problem {
	csvReader := csv.NewReader(reader)
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Failed while extracting lines.")
		fmt.Println("\t", err)
		os.Exit(1)
	}
	var problems []problem
	for _, line := range lines {
		if line[0] == "" || line[1] == "" {
			continue
		}
		problems = append(problems, problem{
			question: line[0],
			answer:   line[1],
		})
	}
	return problems
}

func shuffleProblems(problems []problem) []problem {
	dst := make([]problem, len(problems))
	copy(dst, problems)
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
	return dst
}

type problem struct {
	question string
	answer   string
}

type args struct {
}
