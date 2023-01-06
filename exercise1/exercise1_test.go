package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_parseLinesSuccessfully(t *testing.T) {
	t.Log("Should create a list of problems")
	problems := parseLines(strings.NewReader("a,b\nc,d\n"))

	p1 := problem{
		question: "a",
		answer:   "b",
	}
	p2 := problem{
		question: "c",
		answer:   "d",
	}
	assert.Equal(t, 2, len(problems))
	assert.Equal(t, p1, problems[0])
	assert.Equal(t, p2, problems[1])
}

func Test_parseLinesWithShortLine(t *testing.T) {
	t.Log("Should create a list of problems")
	problems := parseLines(strings.NewReader("a,\nc,d\n"))

	p1 := problem{
		question: "c",
		answer:   "d",
	}
	assert.Equal(t, 1, len(problems))
	assert.Equal(t, p1, problems[0])
}

func Test_shuffleProblems(t *testing.T) {
	t.Log("Returned problems should have the same length")
	problems := []problem{
		{question: "a", answer: "b"},
		{question: "c", answer: "d"},
		{question: "e", answer: "f"},
	}

	result := shuffleProblems(problems)
	assert.Equal(t, len(result), len(problems))
}

func Test_shuffleProblemsCreatesDifferentOrders(t *testing.T) {
	t.Log("Successive shuffles should return different orders")
	problems := []problem{
		{question: "a", answer: "b"},
		{question: "c", answer: "d"},
		{question: "e", answer: "f"},
		{question: "g", answer: "h"},
		{question: "i", answer: "j"},
		{question: "k", answer: "l"},
		{question: "m", answer: "n"},
	}

	r1 := shuffleProblems(problems)
	r2 := shuffleProblems(problems)

	assert.NotEqual(t, r1, r2)
}
