package main

import (
	"testing"
)

func TestInstall(t *testing.T) {
	value := 1
	expected := 2
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}

func TestParseInstallUrl(t *testing.T) {
	value, e := parseInstallURL("user/app")
	expected := "https://github.com/user/app.git"
	if e != nil {
		t.Fatalf("Expected %v, but errored %v:", expected, e)
		return
	}
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}

}
