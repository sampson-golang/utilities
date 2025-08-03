package merge_test

import (
	"testing"

	"github.com/sampson-golang/utilities/container/merge"
)

type TestStruct struct {
	Name    string
	Age     int
	Email   string
	Active  bool
	Height  float64
	Pointer *string
}

func TestStructs_EmptyDestination(t *testing.T) {
	dest := &TestStruct{}
	src := &TestStruct{
		Name:   "John",
		Age:    30,
		Email:  "john@example.com",
		Active: true,
		Height: 5.9,
	}

	merge.Structs(dest, src)

	if dest.Name != "John" {
		t.Errorf("Expected Name to be 'John', got '%s'", dest.Name)
	}
	if dest.Age != 30 {
		t.Errorf("Expected Age to be 30, got %d", dest.Age)
	}
	if dest.Email != "john@example.com" {
		t.Errorf("Expected Email to be 'john@example.com', got '%s'", dest.Email)
	}
	if dest.Active != true {
		t.Errorf("Expected Active to be true, got %t", dest.Active)
	}
	if dest.Height != 5.9 {
		t.Errorf("Expected Height to be 5.9, got %f", dest.Height)
	}
}

func TestStructs_PartiallyFilledDestination(t *testing.T) {
	dest := &TestStruct{
		Name: "Jane",
		Age:  25,
	}
	src := &TestStruct{
		Name:   "John",
		Age:    30,
		Email:  "john@example.com",
		Active: true,
		Height: 5.9,
	}

	merge.Structs(dest, src)

	// All non-empty fields from src should overwrite dest fields
	if dest.Name != "John" {
		t.Errorf("Expected Name to be overwritten to 'John', got '%s'", dest.Name)
	}
	if dest.Age != 30 {
		t.Errorf("Expected Age to be overwritten to 30, got %d", dest.Age)
	}
	// Email, Active, and Height should be set from src
	if dest.Email != "john@example.com" {
		t.Errorf("Expected Email to be 'john@example.com', got '%s'", dest.Email)
	}
	if dest.Active != true {
		t.Errorf("Expected Active to be true, got %t", dest.Active)
	}
	if dest.Height != 5.9 {
		t.Errorf("Expected Height to be 5.9, got %f", dest.Height)
	}
}

func TestStructs_MultipleSources(t *testing.T) {
	dest := &TestStruct{}
	src1 := &TestStruct{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}
	src2 := &TestStruct{
		Active: true,
		Height: 5.9,
		Name:   "Jane",
	}

	merge.Structs(dest, src1, src2)

	if dest.Name != "Jane" {
		t.Errorf("Expected Name to be 'Jane', got '%s'", dest.Name)
	}
	if dest.Age != 30 {
		t.Errorf("Expected Age to be 30, got %d", dest.Age)
	}
	if dest.Email != "john@example.com" {
		t.Errorf("Expected Email to be 'john@example.com', got '%s'", dest.Email)
	}
	if dest.Active != true {
		t.Errorf("Expected Active to be true, got %t", dest.Active)
	}
	if dest.Height != 5.9 {
		t.Errorf("Expected Height to be 5.9, got %f", dest.Height)
	}
}

func TestStructs_WithPointers(t *testing.T) {
	dest := &TestStruct{}
	value := "test pointer"
	src := &TestStruct{
		Name:    "John",
		Pointer: &value,
	}

	merge.Structs(dest, src)

	if dest.Name != "John" {
		t.Errorf("Expected Name to be 'John', got '%s'", dest.Name)
	}
	if dest.Pointer == nil {
		t.Error("Expected Pointer to be set, got nil")
	} else if *dest.Pointer != "test pointer" {
		t.Errorf("Expected Pointer value to be 'test pointer', got '%s'", *dest.Pointer)
	}
}

func TestStructs_ZeroValueShouldNotOverwrite(t *testing.T) {
	dest := &TestStruct{
		Name:   "Jane",
		Age:    25,
		Active: true,
	}
	src := &TestStruct{
		// All fields are zero values
		Name:   "",
		Age:    0,
		Active: false,
	}

	merge.Structs(dest, src)

	// dest should remain unchanged since src has zero values
	if dest.Name != "Jane" {
		t.Errorf("Expected Name to remain 'Jane', got '%s'", dest.Name)
	}
	if dest.Age != 25 {
		t.Errorf("Expected Age to remain 25, got %d", dest.Age)
	}
	if dest.Active != true {
		t.Errorf("Expected Active to remain true, got %t", dest.Active)
	}
}

func TestStructs_BooleanFalseIsZeroValue(t *testing.T) {
	dest := &TestStruct{
		Active: true,
	}
	src := &TestStruct{
		Active: false, // false is zero value for bool
	}

	merge.Structs(dest, src)

	// Active should remain true since false is zero value and shouldn't overwrite
	if dest.Active != true {
		t.Errorf("Expected Active to remain true, got %t", dest.Active)
	}
}
