package utils

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
)

var maxProcs int64

type stop struct {
	error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SetMaxProcs limits the number of concurrent processing goroutines to the given value.
// A value <= 0 clears the limit.
func SetMaxProcs(value int) {
	atomic.StoreInt64(&maxProcs, int64(value))
}

// parallel processes the data in separate goroutines.
func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	limit := int(atomic.LoadInt64(&maxProcs))
	if procs > limit && limit > 0 {
		procs = limit
	}
	if procs > count {
		procs = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}

func getExecDir() (cwd string) {
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

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}
		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			log.Printf("Retry in %v", sleep)

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}

func waitNetworkAvailable() (ok bool) {
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

func setLogger(logFile string) *os.File {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Cannot set logger: %s\n", err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.Println("Started")

	return file
}

func stopLogger(file *os.File) {
	log.Println("Stopped")
	file.Close()
}

func sigIntHook(f func()) {
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		f()
		os.Exit(0)
	}()
}

func createDir(dir string) (err error) {
	_, err = os.Stat(dir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}

	if err != nil {
		log.Fatal(err)
	}

	return
}

func removeContents(dir string) {
	files, err := filepath.Glob(dir)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			log.Println(err)
		}
	}
}

func countFiles(dir string) int {
	files, err := filepath.Glob(dir)
	if err != nil {
		log.Println(err)
		return 0
	}
	return len(files)
}

func splitCommand(cmd string) []string {
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

func executeCommand(command string) error {
	parts := splitCommand(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	return err
}

func getCaptureCommand(cfg Config, imageName string) (string, error) {
	data := struct {
		Image string
		ConfigCamera
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
