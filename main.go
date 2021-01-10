package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	cfg config
)

const (
	// AppVersion is application version
	AppVersion = "v1.0.0"
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

func loadConfiguration() {
	dir := getExecDir()
	if err := os.Chdir(filepath.Dir(dir)); err != nil {
		log.Fatal(err)
	}
	loadConfig(&cfg, "config.yml")
}

func createBaseImageDir() {
	err := createDir(cfg.Settings.BaseImageDir)
	if err != nil {
		log.Printf("Cannot create base image directory: %s", err)
	}
}

func main() {
	var err error

	parseFlags()
	loadConfiguration()

	logFile := setLogger(cfg.Settings.LogFile)
	defer stopLogger(logFile)
	sigIntHook(func() { stopLogger(logFile) })

	waitNetworkAvailable()
	createBaseImageDir()
	sendEmail("CAM START", "Camera started", nil, cfg.SMTP)

	// main loop
	var lastWeekdayHour string
	var lastImage string
	lastAlert := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	for {
		// update time
		currentTime := time.Now()
		weekday := fmt.Sprintf("%02d", currentTime.Weekday())
		weekdayHour := fmt.Sprintf("%s%02d", weekday, currentTime.Hour())

		// update directory
		imagePath := filepath.Join(cfg.Settings.BaseImageDir, weekday)
		allImagesMask := filepath.Join(imagePath, weekdayHour+"-*.jpg")
		err = createDir(imagePath)
		if err != nil {
			log.Printf("Cannot create directory: %s\n", err)
		}
		if weekdayHour != lastWeekdayHour {
			if len(lastWeekdayHour) > 0 {
				removeContents(allImagesMask)
			}
			lastWeekdayHour = weekdayHour
		}

		// get new image name
		imageName := filepath.Join(imagePath, fmt.Sprintf("%s-%04d.jpg", weekdayHour, 1+countFiles(allImagesMask)))
		var sidx float32
		dirty := false

		captureCommand, err := getCaptureCommand(imageName)
		if err != nil {
			log.Println(err)
			log.Fatal("Failed to create capture command\n")
		}
		err = retry(5, 1*time.Second, func() error { return executeCommand(captureCommand) })
		if err != nil {
			log.Printf("Failed to capture image: %s\n", err)
			goto nextIteration
		}

		if len(lastImage) > 0 {
			sidx, err = imageSimilarityIndexFile(imageName, lastImage, cfg.Settings.Sensitivity)
			if err != nil {
				log.Printf("Failed to calculate similarity index: %s\n", err)
				goto nextIteration
			}

			fmt.Printf("Similarity index = %.2f\n", sidx)

			if sidx < cfg.Settings.KeepThreshold {
				err = os.Remove(imageName)
				if err != nil {
					log.Println(err)
				}
			}

			// upload to FTP
			if sidx > cfg.Settings.UploadThreshold {
				err = uploadFTP(cfg.FTP, imageName, weekday)
				if err != nil {
					log.Printf("Failed to upload to FTP: %s\n", err)
				}
				dirty = true
			}

			// send email alert
			if sidx > cfg.Settings.EmailThreshold && currentTime.Sub(lastAlert).Seconds() > 300.0 {
				err = sendEmail("CAMERA ALERT", "", []string{imageName}, cfg.SMTP)
				if err != nil {
					log.Printf("Failed to send email alert: %s\n", err)
				}
				dirty = true
				lastAlert = time.Now()
			}
		}

		if dirty || len(lastImage) == 0 {
			lastImage = imageName
		}

	nextIteration:

		// pause if necessary, we do not want to overload the loop in case of issues
		endTime := time.Now()
		elapsed := endTime.Sub(currentTime).Milliseconds()
		sleepTime := time.Duration(1000 - elapsed)
		if sleepTime > 0 {
			time.Sleep(sleepTime * time.Millisecond)
		}

		fmt.Printf("Elapsed time %0.2fs, last image=%s\n", float64(elapsed)/1000, imageName)
	}
}
