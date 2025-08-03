package output_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/sampson-golang/utilities/output"
)

func executeTest(t *testing.T, run func(), expected string) {
	oldStdout := os.Stdout
	var buf bytes.Buffer

	r, w, _ := os.Pipe()
	defer func() {
		os.Stdout = oldStdout
		w.Close()
	}()

	os.Stdout = w
	run()
	w.Close()
	os.Stdout = oldStdout

	io.Copy(&buf, r)

	result := buf.String()

	if result != expected {
		t.Errorf("Expected output to be %q, got %q", expected, result)
	}
}

func TestPrettyPrint_OutputsToStdout(t *testing.T) {
	input := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	executeTest(
		t,
		func() {
			output.PrettyPrint(input)
		},
		"{\n  \"age\": 30,\n  \"name\": \"John\"\n}\n",
	)

	t.Run("with custom indent", func(t *testing.T) {
		executeTest(
			t,
			func() {
				output.PrettyPrint(input, "\t")
			},
			"{\n\t\"age\": 30,\n\t\"name\": \"John\"\n}\n",
		)
	})
}

func TestPrettyPrint_SimpleString(t *testing.T) {
	executeTest(
		t,
		func() {
			output.PrettyPrint("hello world")
		},
		"\"hello world\"\n",
	)

	t.Run("with custom indent", func(t *testing.T) {
		executeTest(
			t,
			func() {
				output.PrettyPrint("hello world", "\t")
			},
			"\"hello world\"\n",
		)
	})
}

func TestPrettyPrint_Array(t *testing.T) {
	input := []string{"apple", "banana", "cherry"}
	executeTest(
		t,
		func() {
			output.PrettyPrint(input)
		},
		"[\n  \"apple\",\n  \"banana\",\n  \"cherry\"\n]\n",
	)

	t.Run("with custom indent", func(t *testing.T) {
		executeTest(
			t,
			func() {
				output.PrettyPrint(input, "\t")
			},
			"[\n\t\"apple\",\n\t\"banana\",\n\t\"cherry\"\n]\n",
		)
	})
}

func TestPrettyPrint_NilInput(t *testing.T) {
	executeTest(
		t,
		func() {
			output.PrettyPrint(nil)
		},
		"null\n",
	)

	t.Run("with custom indent", func(t *testing.T) {
		executeTest(
			t,
			func() {
				output.PrettyPrint(nil, "\t")
			},
			"null\n",
		)
	})
}

func TestPrettyPrint_EmptyMap(t *testing.T) {
	executeTest(
		t,
		func() {
			output.PrettyPrint(map[string]interface{}{})
		},
		"{}\n",
	)

	t.Run("with custom indent", func(t *testing.T) {
		executeTest(
			t,
			func() {
				output.PrettyPrint(map[string]interface{}{}, "\t")
			},
			"{}\n",
		)
	})
}
