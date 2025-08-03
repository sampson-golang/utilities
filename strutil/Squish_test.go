package strutil_test

import (
	"testing"

	"github.com/sampson-golang/utilities/strutil"
)

func TestSquish_BasicWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"leading spaces", "   hello", "hello"},
		{"trailing spaces", "hello   ", "hello"},
		{"both leading and trailing", "   hello   ", "hello"},
		{"internal double spaces", "hello  world", "hello world"},
		{"internal triple spaces", "hello   world", "hello world"},
		{"multiple internal spaces", "hello    world    test", "hello world test"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := strutil.Squish(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSquish_DifferentWhitespaceTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"tabs", "hello\t\tworld", "hello world"},
		{"newlines", "hello\n\nworld", "hello world"},
		{"mixed whitespace", "hello \t\n world", "hello world"},
		{"leading tab", "\thello", "hello"},
		{"trailing newline", "hello\n", "hello"},
		{"complex mixed", "\t\n  hello \t\n world \n\t ", "hello world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := strutil.Squish(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSquish_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
		{"only tabs", "\t\t\t", ""},
		{"only newlines", "\n\n\n", ""},
		{"only mixed whitespace", " \t\n ", ""},
		{"single character", "a", "a"},
		{"single character with spaces", " a ", "a"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := strutil.Squish(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSquish_NoWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "hello", "hello"},
		{"two words", "hello world", "hello world"},
		{"multiple words", "hello world test", "hello world test"},
		{"punctuation", "hello,world!", "hello,world!"},
		{"numbers", "123456", "123456"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := strutil.Squish(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSquish_PreservesNonWhitespace(t *testing.T) {
	input := "  hello,   world!  \t\nHow  are\n\tyou?  "
	expected := "hello, world! How are you?"

	result := strutil.Squish(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSquish_MultilineText(t *testing.T) {
	input := `  This is
	a multiline
		string with
	various   spacing  `

	expected := "This is a multiline string with various spacing"

	result := strutil.Squish(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
