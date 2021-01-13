package process

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kornelkabele/watchdog/internal/cfg"
	"github.com/kornelkabele/watchdog/internal/email"
	"github.com/kornelkabele/watchdog/internal/file"
	ftp "github.com/kornelkabele/watchdog/internal/ftp"
	img "github.com/kornelkabele/watchdog/internal/image"
	"github.com/kornelkabele/watchdog/internal/system"
)

var (
	lastWeekdayHour string
	lastImage       string
	lastAlert       time.Time
)

func init() {
	lastAlert = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
}

// Process runs full update-capture-compare-store-upload-email iteration
func Process() {
	// update time
	currentTime := time.Now()
	weekday := fmt.Sprintf("%02d", currentTime.Weekday())
	weekdayHour := fmt.Sprintf("%s%02d", weekday, currentTime.Hour())

	// update directory
	imagePath := filepath.Join(cfg.Settings.BaseImageDir, weekday)
	allImagesMask := filepath.Join(imagePath, weekdayHour+"-*.jpg")
	err := file.CreateDir(imagePath)
	if err != nil {
		log.Fatalf("Cannot create directory: %s\n", err)
	}
	if weekdayHour != lastWeekdayHour {
		if len(lastWeekdayHour) > 0 {
			file.RemoveContents(allImagesMask)
		}
		lastWeekdayHour = weekdayHour
	}

	// capture new image
	numFiles, err := file.CountFiles(allImagesMask)
	if err != nil {
		log.Printf("Failed to count number of files in directory: %s\n", err)
	}
	imageName := filepath.Join(imagePath, fmt.Sprintf("%s-%04d.jpg", weekdayHour, 1+numFiles))
	captureCommand, err := system.GetCaptureCommand(imageName)
	if err != nil {
		log.Printf("Failed to create capture command: %s\n", err)
	}
	err = retry(5, 1*time.Second, func() error { return system.ExecuteCommand(captureCommand, 10*time.Second) })
	if err != nil {
		log.Printf("Failed to capture image: %s\n", err)
		err = email.SendEmail("CAMERA CAPTURE FAILURE", fmt.Sprintf("Failed to capture camera: %s", err), nil)
		if err != nil {
			log.Printf("Failed to send email failure: %s\n", err)
		}
		return
	}

	// keep if there is no reference
	if len(lastImage) == 0 {
		lastImage = imageName
		return
	}

	// compute similarity index
	sidx, err := img.ImageSimilarityIndexFile(imageName, lastImage, cfg.Settings.Sensitivity)
	if err != nil {
		log.Printf("Failed to calculate similarity index: %s\n", err)
		return
	}

	fmt.Printf("Similarity index = %.2f (%s)\n", sidx, imageName)

	// remove from local directory if too similar
	if sidx < cfg.Settings.KeepThreshold {
		err = os.Remove(imageName)
		if err != nil {
			log.Println(err)
		}
	}

	// upload to FTP
	if sidx > cfg.Settings.UploadThreshold {
		err = ftp.UploadFTP(imageName, weekday)
		if err != nil {
			log.Printf("Failed to upload to FTP: %s\n", err)
			err = email.SendEmail("CAMERA FTP FAILURE", fmt.Sprintf("Failed to upload to FTP: %s\n", err), nil)
			if err != nil {
				log.Printf("Failed to send FTP failure: %s\n", err)
			}
		}
		lastImage = imageName
	}

	// send email alert
	if sidx > cfg.Settings.EmailThreshold && currentTime.Sub(lastAlert).Seconds() > float64(cfg.Settings.EmailInterval) {
		err = email.SendEmail("CAMERA ALERT", "", []string{imageName})
		if err != nil {
			log.Printf("Failed to send email alert: %s\n", err)
		}
		lastImage = imageName
		lastAlert = time.Now()
	}
}
