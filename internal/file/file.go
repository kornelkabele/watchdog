package file

import (
	"log"
	"os"
	"path/filepath"
)

// CreateDir creates a directory if it does not exist
func CreateDir(dir string) (err error) {
	_, err = os.Stat(dir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}

	return
}

// RemoveContents removes directory contents
func RemoveContents(dir string) {
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

// CountFiles counts number of files in directory
func CountFiles(dir string) (int, error) {
	files, err := filepath.Glob(dir)
	return len(files), err
}
