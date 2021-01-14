package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kornelkabele/watchdog/internal/cfg"
	"github.com/kornelkabele/watchdog/internal/email"
	"github.com/kornelkabele/watchdog/internal/file"
	"github.com/kornelkabele/watchdog/internal/logger"
	"github.com/kornelkabele/watchdog/internal/process"
	"github.com/kornelkabele/watchdog/internal/system"
)

var (
	config cfg.Config
)

const (
	// AppVersion is application version
	AppVersion = "v1.0.0"
	// ConfigFile is default configuration file
	ConfigFile = "config.yml"
)

func parseFlags() {
	// parse flags
	version := flag.Bool("v", false, "prints current application version")
	flag.Parse()
	if *version {
		fmt.Println("Version:\t", AppVersion)
		os.Exit(0)
	}
}

func createBaseImageDir() {
	err := file.CreateDir(cfg.Settings.BaseImageDir)
	if err != nil {
		log.Fatalf("Cannot create base image directory: %s", err)
	}
}

func main() {
	var err error

	dir := system.GetExecDir()
	if err := os.Chdir(filepath.Dir(dir)); err != nil {
		log.Fatal(err)
	}

	parseFlags()
	cfg.LoadConfig(ConfigFile)

	logFile, err := logger.SetLogger(cfg.Settings.LogFile)
	if err != nil {
		log.Fatalf("Cannot set logger: %s\n", err)
	}
	defer logger.StopLogger(logFile)
	system.SigIntHook(func() { logger.StopLogger(logFile) })

	system.WaitNetworkAvailable()
	createBaseImageDir()
	email.SendEmail(fmt.Sprintf("CAMERA START: %s", cfg.Settings.Id),
		fmt.Sprintf("%s Camera started", time.Now().Format(time.RFC3339)),
		nil)

	// main loop
	for {
		// update time
		currentTime := time.Now()

		process.Process()

		// pause if necessary, we do not want to overload the loop in case of issues
		endTime := time.Now()
		elapsed := endTime.Sub(currentTime).Milliseconds()
		sleepTime := time.Duration(1000 - elapsed)
		if sleepTime > 0 {
			time.Sleep(sleepTime * time.Millisecond)
		}

		fmt.Printf("Elapsed time %0.2fs\n", float64(elapsed)/1000)
	}
}
