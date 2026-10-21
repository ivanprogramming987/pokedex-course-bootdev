package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	// initalize test cases
	test_cases := []struct{
		input    string
		expected []string
		id       int
	}{
	{
		input:    "!test case? wi&th p(unctuatio)n *marks*",
		expected: []string{"!test", "case?", "wi&th", "p(unctuatio)n", "*marks*"},
		id:       1,
	},
	{
		input:    "         test case with outside whitespace         ",
		expected: []string{"test", "case", "with", "outside", "whitespace"},
		id:       2,
	},
	{
		input:    "test      case        with      inside       whitespace",
		expected: []string{"test", "case", "with", "inside", "whitespace"},
		id:       3,
	},
	{
		input:    "Test casE wITh CAPITAL lEtteRs",
		expected: []string{"test", "case", "with", "capital", "letters"},
		id:       4,
	},
	}

	for _, c := range test_cases {
		actual := cleanInput(c.input)

		// check the length of the slice returned
		if len(c.expected) != len(actual) {
			t.Errorf("TEST CASE %d FAILED: \nInput: %s\nExpected: %v\nActual: %v", c.id, c.input, c.expected, actual)
			continue
		}

		// check that every item matched
		for i, _ := range c.expected {
			if c.expected[i] != actual[i] {
				t.Errorf("TEST CASE %d FAILED: \nInput: %s\nExpected: %v\nActual: %v", c.id, c.input, c.expected, actual)
			}
		}
	}
}
