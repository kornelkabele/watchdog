package utils

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// Stop if returned as error stops Retry function
type Stop struct {
	error
}

// Retry tries to execute function f specified number of times
func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(Stop); ok {
			// Return the original error for later checking
			return s.error
		}
		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			log.Printf("Retry in %v", sleep)

			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}

// WaitNetworkAvailable waits for network availability, important if script is started from cron using wifi connection
func WaitNetworkAvailable() (ok bool) {
	err := Retry(5, 5*time.Second, func() (err error) {
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
func ExecuteCommand(command string) error {
	parts := SplitCommand(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	err := cmd.Run()
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
