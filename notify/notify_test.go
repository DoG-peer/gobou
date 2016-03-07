package notify

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	args := []string{"notify-send", "hoge"}
	cmd := exec.Command(args[0], args[1:]...)
	value := "notify-send hoge"
	expected := strings.Join(cmd.Args, " ")
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}

func TestNotifyArgs(t *testing.T) {
	args := NotifySendArgs{
		text: "hoge",
	}
	cmd := NotifySend(args)
	value := strings.Join(cmd.Args, " ")
	expected := "notify-send hoge"
	if value != expected {
		t.Fatalf("Expected %v, but %v:", expected, value)
	}
}
