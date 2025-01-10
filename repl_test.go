package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "PIKACHU  Charmander    Bulbasaur",
			expected: []string{"pikachu", "charmander", "bulbasaur"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "Testing   Multiple   Spaces",
			expected: []string{"testing", "multiple", "spaces"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check slice lengths
		if len(actual) != len(c.expected) {
			t.Errorf("Length mismatch for input '%s': expected %d, got %d",
				c.input, len(c.expected), len(actual))
			continue
		}
		// Check each word
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word mismatch for input '%s' at position %d: expected '%s', got '%s'",
					c.input, i, expectedWord, word)
			}
		}
	}
}
