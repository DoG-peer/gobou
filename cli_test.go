package main

import (
	"testing"
)

func TestCli(t *testing.T) {
	value := 1
	expected := 2
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}
