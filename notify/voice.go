package notify

import (
	"fmt"
	"github.com/DoG-peer/gobou/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Jtalk is about open_jtalk
type Jtalk struct {
	Voice      string
	Dictionary string
	TmpDir     string
}

// Play is
func (j *Jtalk) Play(m gobou.JtalkMessage) {
	if m.IsNone() {
		return
	}
	path := filepath.Join(j.TmpDir, "tmp.wav")
	WaitForWrite(path)
	exec.Command("aplay", "--quiet", path).Run()
	os.Remove(path)
}

// MakeWavFile makes .wav file
func (j *Jtalk) MakeWavFile(words string) error {
	path := filepath.Join(j.TmpDir, "tmp.wav")
	os.Remove(path)
	cmd := exec.Command("open_jtalk",
		"-x", j.Dictionary,
		"-m", j.Voice,
		"-ow", path,
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	if err = cmd.Start(); err != nil {
		return err
	}
	io.WriteString(stdin, words)
	return nil
}

// WaitForWrite waits until the file has been written
func WaitForWrite(path string) error {
	maxStep := 1000
	fs, e := os.Stat(path)
	size := int64(0)
	for i := 0; i < maxStep; i++ {
		if e == nil {
			if size != 0 && fs.Size() == size {
				return nil
			}
			size = fs.Size()
		}
		time.Sleep(10 * time.Millisecond)
		fs, e = os.Stat(path)
	}
	return fmt.Errorf("TimeOut for %s", path)
}
