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

func createImageDir() {
	err := file.CreateDir(cfg.Settings.ImageDir)
	if err != nil {
		log.Fatalf("Cannot create image directory: %s", err)
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
	fmt.Printf("Settings:\n")
	fmt.Printf("Id: %s\n", cfg.Settings.Id)
	fmt.Printf("Camera: %s:%d\n", cfg.Camera.Host, cfg.Camera.Port)
	fmt.Printf("Email: %s:%d\n", cfg.SMTP.Host, cfg.SMTP.Port)
	fmt.Printf("FTP: %s:%d\n", cfg.FTP.Host, cfg.FTP.Port)
	fmt.Printf("Image dir: %s\n", cfg.Settings.ImageDir)
	fmt.Printf("Log file: %s\n", cfg.Settings.LogFile)

	logger, err := logger.New(cfg.Settings.LogFile)
	if err != nil {
		log.Fatalf("Cannot set logger: %s\n", err)
	}
	defer logger.Close()
	system.SigIntHook(func() { logger.Close() })

	system.WaitNetworkAvailable()
	createImageDir()
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
