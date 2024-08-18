package main

import (
	"testing"
)

func TestCalculateSteps(t *testing.T) {
	tests := []struct {
		password string
		steps    int
	}{
		{"aA1", 3},
		{"1445D1cd", 0},
		{"", 6},
		{"abaAcA161", 3},
		{"abcabcabcabcabcabcabcabcabcabc", 12},
	}

	for _, test := range tests {
		result := calculateSteps(test.password)
		if result != test.steps {
			t.Errorf("For password %s, expected %d steps, but got %d", test.password, test.steps, result)
		}
	}
}
