package system

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/kornelkabele/watchdog/internal/cfg"
)

// GetExecDir returns executable directory
func GetExecDir() (cwd string) {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cwd, err = filepath.EvalSymlinks(path)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// WaitNetworkAvailable waits for network availability, important if script is started from cron using wifi connection
func WaitNetworkAvailable() (ok bool) {
	err := retry(5, 5*time.Second, func() (err error) {
		_, err = http.Get("http://clients1.google.com/generate_204")
		return
	})
	if err != nil {
		log.Printf("Network not available: %s\n", err)
		return false
	}
	return true
}

// SigIntHook attaches function to ^C interrupt signal
func SigIntHook(f func()) {
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		f()
		os.Exit(0)
	}()
}

// SplitCommand splits shell command to string slice
func SplitCommand(cmd string) []string {
	lastQuote := rune(0)
	return strings.FieldsFunc(cmd, func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return true
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return true
		default:
			return unicode.IsSpace(c)
		}
	})
}

// ExecuteCommand executes shell command
func ExecuteCommand(command string, duration time.Duration) error {
	parts := SplitCommand(command)

	var cmd *exec.Cmd
	if duration == 0 {
		cmd = exec.Command(parts[0], parts[1:]...)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
	}

	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		log.Println(string(out))
	}
	if err != nil {
		log.Println(err)
	}

	return err
}

// GetCaptureCommand processes ffmpeg template and returns shell command
func GetCaptureCommand(imageName string) (string, error) {
	data := struct {
		Image string
		cfg.ConfigCamera
	}{
		imageName,
		cfg.Camera,
	}
	tmpl, err := template.New("Action").Parse(cfg.Settings.FFmpegCmd)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
