package notify

import (
	"os/exec"
)

type Notifier interface {
	NotifyCommand() exec.Cmd
}

// NotifySendArgs is argment for notify-send
// urgency must be low, normal, or critical
// expireTime is millisecond
// icon is
type NotifySendArgs struct {
	urgency    string
	expireTime string
	appName    string
	icon       string
	category   string
	text       string
	title      string
}

func Notify(s string) error {
	return NotifySend(NotifySendArgs{text: s}).Run()
}

func (args NotifySendArgs) Args() []string {
	inputs := []string{}
	if args.urgency != "" {
		inputs = append(inputs, "-t", args.urgency)
	}
	if args.expireTime != "" {
		inputs = append(inputs, "-t", args.expireTime)
	}
	if args.appName != "" {
		inputs = append(inputs, "-a", args.appName)
	}
	if args.category != "" {
		inputs = append(inputs, "-c", args.category)
	}
	if args.icon != "" {
		inputs = append(inputs, "-i", args.icon)
	}
	if args.title != "" {
		inputs = append(inputs, args.title)
	}
	if args.text != "" {
		inputs = append(inputs, args.text)
	}
	return inputs
}

func NotifySend(args NotifySendArgs) *exec.Cmd {

	return exec.Command("notify-send", args.Args()...)

}
