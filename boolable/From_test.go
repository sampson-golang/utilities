package boolable_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sampson-golang/utilities/boolable"
)

func TestBoolable(t *testing.T) {
	falseStrings := []string{
		"",
		"0",
		"f",
		"false",
		"off",
		"n",
		"no",
	}

	bools := []bool{true, false}

	falseValues := []interface{}{false, 0, nil}
	for _, s := range falseStrings {
		falseValues = append(falseValues, s)
	}

	t.Run("bool", func(t *testing.T) {
		for _, value := range bools {
			t.Run(fmt.Sprintf("%v", value), func(t *testing.T) {
				result := boolable.From(value)
				if result != value {
					t.Errorf("Boolable(%v) = %v, expected %v", value, result, value)
				}
			})
		}
	})

	t.Run("int", func(t *testing.T) {
		t.Run("0 is false", func(t *testing.T) {
			result := boolable.From(0)
			if result != false {
				t.Errorf("Boolable(0) = %v, expected false", result)
			}
		})

		t.Run("positive int is true", func(t *testing.T) {
			for i := 1; i < 101; i++ {
				result := boolable.From(i)
				if result != true {
					t.Errorf("Boolable(%v) = %v, expected true", i, result)
				}
			}
		})

		t.Run("negative int is true", func(t *testing.T) {
			for i := -1; i > -101; i-- {
				result := boolable.From(i)
				if result != true {
					t.Errorf("Boolable(%v) = %v, expected true", i, result)
				}
			}
		})
	})

	t.Run("string", func(t *testing.T) {
		t.Run("false values", func(t *testing.T) {
			t.Run("empty", func(t *testing.T) {
				result := boolable.From("")
				if result != false {
					t.Errorf("Boolable(\"\") = %v, expected false", result)
				}
			})

			for _, value := range falseStrings {
				if value == "" {
					continue
				}

				t.Run(value, func(t *testing.T) {
					result := boolable.From(value)
					if result != false {
						t.Errorf("Boolable(%v) = %v, expected false", value, result)
					}
				})
			}

			t.Run("are not case sensitive", func(t *testing.T) {
				for _, value := range falseStrings {
					lower := strings.ToLower(value)
					upper := strings.ToUpper(value)

					result := boolable.From(lower)
					if result != false {
						t.Errorf("Boolable(%v) = %v, expected false", lower, result)
					}

					result = boolable.From(upper)
					if result != false {
						t.Errorf("Boolable(%v) = %v, expected false", upper, result)
					}

					for i := 0; i < len(value); i++ {
						mid := value[:i] + strings.ToUpper(string(value[i])) + value[i+1:]
						result := boolable.From(mid)
						if result != false {
							t.Errorf("Boolable(%v) = %v, expected false", mid, result)
						}
					}
				}
			})

			t.Run("must be exact except for casing", func(t *testing.T) {
				for _, value := range falseStrings {
					result := boolable.From(value + " ")
					if result != true {
						t.Errorf("Boolable(%v) = %v, expected true", value, result)
					}

					result = boolable.From(" " + value)
					if result != true {
						t.Errorf("Boolable(%v) = %v, expected true", value, result)
					}

					if value == "" {
						continue
					}

					for _, dup := range falseStrings {
						if dup == "" {
							continue
						}

						result = boolable.From(value + dup)
						if result != true {
							t.Errorf("Boolable(%v) = %v, expected true", value+dup, result)
						}
					}
				}
			})
		})

		t.Run("everything else is true", func(t *testing.T) {
			for _, value := range []string{" ", "true", "t", "1", "00", "`"} {
				result := boolable.From(value + "a")
				if result != true {
					t.Errorf("Boolable(%v) = %v, expected true", value, result)
				}
			}
		})
	})

	t.Run("nil is false", func(t *testing.T) {
		result := boolable.From(nil)
		if result != false {
			t.Errorf("Boolable(nil) = %v, expected false", result)
		}

		var nilInterface interface{}
		result = boolable.From(nilInterface)
		if result != false {
			t.Errorf("Boolable(nilInterface) = %v, expected false", result)
		}

		var nilStruct *struct{} = nil
		result = boolable.From(nilStruct)
		if result != false {
			t.Errorf("Boolable(nilStruct) = %v, expected false", result)
		}

		var nilSlice *[]bool = nil
		result = boolable.From(nilSlice)
		if result != false {
			t.Errorf("Boolable(%v) = %v, expected false", nilSlice, result)
		}

		var nilMap *map[string]bool = nil
		result = boolable.From(nilMap)
		if result != false {
			t.Errorf("Boolable(%v) = %v, expected false", nilMap, result)
		}
	})

	t.Run("struct", func(t *testing.T) {
		result := boolable.From(struct{}{})
		if result != true {
			t.Errorf("Boolable(struct{}) = %v, expected true", result)
		}
	})

	t.Run("slice", func(t *testing.T) {
		result := boolable.From([]bool{true})
		if result != true {
			t.Errorf("Boolable([]bool) = %v, expected true", result)
		}

		result = boolable.From([]bool{false})
		if result != true {
			t.Errorf("Boolable([]bool) = %v, expected false", result)
		}

		result = boolable.From([]bool{})
		if result != true {
			t.Errorf("Boolable([]bool) = %v, expected false", result)
		}

		result = boolable.From(falseStrings)
		if result != true {
			t.Errorf("Boolable([]string) = %v, expected true", result)
		}
	})

	t.Run("optional dereference parameter", func(t *testing.T) {
		slice := []bool{false}

		t.Run("manages pointer dereferencing", func(t *testing.T) {
			for _, value := range falseStrings {
				result := boolable.From(&value)
				if result != false {
					if value == "" {
						value = "\"\""
					}

					t.Errorf("Boolable(&%v) = %v, expected false", value, result)
				}

				result = boolable.From(&value, true)
				if result != false {
					t.Errorf("Boolable(&%v, true) = %v, expected false", value, result)
				}

				result = boolable.From(&value, false)
				if result != true {
					t.Errorf("Boolable(&%v, false) = %v, expected true", value, result)
				}
			}

			for _, value := range []bool{true, false} {
				result := boolable.From(&value)
				if result != value {
					t.Errorf("Boolable(&%v) = %v, expected %v", value, result, value)
				}

				result = boolable.From(&value, true)
				if result != value {
					t.Errorf("Boolable(&%v, true) = %v, expected %v", value, result, value)
				}

				result = boolable.From(&value, false)
				if result != true {
					t.Errorf("Boolable(&%v, false) = %v, expected true", value, result)
				}
			}

			for i := -101; i < 101; i++ {
				expected := i != 0
				result := boolable.From(&i)
				if result != expected {
					t.Errorf("Boolable(&%v) = %v, expected %v", i, result, expected)
				}

				result = boolable.From(&i, true)
				if result != expected {
					t.Errorf("Boolable(&%v, true) = %v, expected %v", i, result, expected)
				}

				result = boolable.From(&i, false)
				if result != true {
					t.Errorf("Boolable(&%v, false) = %v, expected true", i, result)
				}
			}

			if result := boolable.From(&slice); result != true {
				t.Errorf("Boolable(&slice) = %v, expected true", result)
			}

			if result := boolable.From(&slice, true); result != true {
				t.Errorf("Boolable(&slice, true) = %v, expected true", result)
			}

			if result := boolable.From(&slice, false); result != true {
				t.Errorf("Boolable(&slice, false) = %v, expected true", result)
			}
		})

		t.Run("has no effect on non-pointer values", func(t *testing.T) {
			for _, value := range falseValues {
				result := boolable.From(value)
				if result != false {
					switch v := value.(type) {
					case string:
						if v == "" {
							value = "\"\""
						}
					}

					t.Errorf("Boolable(&%v) = %v, expected false", value, result)
				}
				result = boolable.From(value, true)
				if result != false {
					t.Errorf("Boolable(&%v, true) = %v, expected false", value, result)
				}
				result = boolable.From(value, false)
				if result != false {
					t.Errorf("Boolable(&%v, false) = %v, expected false", value, result)
				}
			}

			for _, value := range []interface{}{
				true,
				42,
				"hello",
				[]bool{true},
				map[string]bool{"key": true},
				struct{}{},
			} {
				result := boolable.From(value)
				if result != true {
					t.Errorf("Boolable(%v) = %v, expected true", value, result)
				}

				result = boolable.From(value, true)
				if result != true {
					t.Errorf("Boolable(%v, true) = %v, expected true", value, result)
				}

				result = boolable.From(value, false)
				if result != true {
					t.Errorf("Boolable(%v, false) = %v, expected true", value, result)
				}
			}
		})
	})
}

// Benchmark tests
func BenchmarkBoolable(b *testing.B) {
	trueVal := true
	falseVal := false
	zeroInt := 0
	nonZeroInt := 42
	emptyStr := ""
	falseStr := "false"
	slice := []bool{false}
	mapVal := map[string]bool{"key": false}
	var nilVal interface{} = nil
	structVal := struct{ Name string }{Name: "test"}

	testCases := []interface{}{
		emptyStr,
		trueVal,
		falseVal,
		zeroInt,
		nonZeroInt,
		falseStr,
		slice,
		mapVal,
		nilVal,
		structVal,
		&emptyStr,
		&trueVal,
		&falseVal,
		&zeroInt,
		&nonZeroInt,
		&falseStr,
		&slice,
		&mapVal,
		&nilVal,
		&structVal,
		(*bool)(nil),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			boolable.From(tc)
			boolable.From(tc, true)
			boolable.From(tc, false)
		}
	}
}
