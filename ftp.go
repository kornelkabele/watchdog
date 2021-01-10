package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/secsy/goftp"
)

func uploadFTP(cfg configFTP, src, dst string) error {
	config := goftp.Config{
		User:               cfg.User,
		Password:           cfg.Pass,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		TLSMode:            goftp.TLSImplicit,
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	client, err := goftp.DialConfig(config, cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Stat(dst)
	if err != nil {
		_, err := client.Mkdir(dst)
		if err != nil {
			return err
		}
	}

	target := "/" + dst + "/" + filepath.Base(src)
	err = client.Store(target, f)
	if err != nil {
		return err
	}
	return nil
}
