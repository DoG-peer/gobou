package gobou

import (
	"github.com/dullgiulio/pingo"
)

// Register connect your type to pingo
func Register(obj interface{}) {
	pingo.Register(obj)
	pingo.Run()
}

// StdoutMessage is for stdout
type StdoutMessage struct {
	Text string
}

// NotifyMessage is for notify-send command
type NotifyMessage struct {
	Title string
	Text  string
	Icon  string
}

// JtalkMessage is for open_jtalk command
type JtalkMessage struct {
	Text string
}

// Message should be used by main plugin tasks
type Message struct {
	StdoutMessage StdoutMessage
	NotifyMessage NotifyMessage
	JtalkMessage  JtalkMessage
}

// IsNone for StdoutMessage
func (m *StdoutMessage) IsNone() bool {
	return m.Text == ""
}

// IsNone for NotifyMessage
func (m *NotifyMessage) IsNone() bool {
	return m.Text == "" && m.Title == ""
}

// IsNone for JtalkMessage
func (m *JtalkMessage) IsNone() bool {
	return m.Text == ""
}

// IsNone for Message
func (m *Message) IsNone() bool {
	return m.StdoutMessage.IsNone() && m.NotifyMessage.IsNone() && m.JtalkMessage.IsNone()
}

// Print returns simple message for stdout
func Print(s string) Message {
	return Message{
		StdoutMessage: StdoutMessage{
			Text: s,
		},
	}
}

// Notify returns simple message for notify-send
func Notify(s string) Message {
	return Message{
		NotifyMessage: NotifyMessage{
			Text: s,
		},
	}
}

// NotifyWithTitle returns simple message for notify-send
func NotifyWithTitle(title, text string) Message {
	return Message{
		NotifyMessage: NotifyMessage{
			Title: title,
			Text:  text,
		},
	}
}

// Say returns simple message for open_jtalk
func Say(s string) Message {
	return Message{
		JtalkMessage: JtalkMessage{
			Text: s,
		},
	}
}
