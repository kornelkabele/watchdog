package main

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type configCamera struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type configFTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type configSMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Sender   string `yaml:"sender"`
	Receiver string `yaml:"receiver"`
}

type configSettings struct {
	Sensitivity     float32 `yaml:"sensitivity"`
	KeepThreshold   float32 `yaml:"keepThreshold"`
	UploadThreshold float32 `yaml:"uploadThreshold"`
	EmailThreshold  float32 `yaml:"emailThreshold"`
	BaseImageDir    string  `yaml:"baseImageDir"`
	LogFile         string  `yaml:"logFile"`
	FFmpegCmd       string  `yaml:"ffmpegCmd"`
}

type config struct {
	Camera   configCamera   `yaml:"camera"`
	FTP      configFTP      `yaml:"ftp"`
	SMTP     configSMTP     `yaml:"smtp"`
	Settings configSettings `yaml:"settings"`
}

func loadConfig(cfg *config, configFile string) {
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	loadEnvSecrets(cfg)
	validateConfig(cfg)
}

func loadEnvSecrets(cfg *config) {
	var err error
	if os.Getenv("CAMERA_HOST") != "" {
		cfg.Camera.Host = os.Getenv("CAMERA_HOST")
	}
	if os.Getenv("CAMERA_PORT") != "" {
		cfg.Camera.Port, err = strconv.Atoi(os.Getenv("CAMERA_PORT"))
		if err != nil {
			log.Fatal(err)
		}
	}
	if os.Getenv("CAMERA_USER") != "" {
		cfg.Camera.User = os.Getenv("CAMERA_USER")
	}
	if os.Getenv("CAMERA_PASS") != "" {
		cfg.Camera.Pass = os.Getenv("CAMERA_PASS")
	}
	if os.Getenv("FTP_HOST") != "" {
		cfg.FTP.Host = os.Getenv("FTP_HOST")
	}
	if os.Getenv("FTP_PORT") != "" {
		cfg.FTP.Port, err = strconv.Atoi(os.Getenv("FTP_PORT"))
		if err != nil {
			log.Fatal(err)
		}
	}
	if os.Getenv("FTP_USER") != "" {
		cfg.FTP.User = os.Getenv("FTP_USER")
	}
	if os.Getenv("FTP_PASS") != "" {
		cfg.FTP.Pass = os.Getenv("FTP_PASS")
	}
	if os.Getenv("SMTP_HOST") != "" {
		cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	}
	if os.Getenv("SMTP_PORT") != "" {
		cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
		if err != nil {
			log.Fatal(err)
		}
	}
	if os.Getenv("SMTP_USER") != "" {
		cfg.SMTP.User = os.Getenv("SMTP_USER")
	}
	if os.Getenv("SMTP_PASS") != "" {
		cfg.SMTP.Pass = os.Getenv("SMTP_PASS")
	}
	if os.Getenv("SMTP_SENDER") != "" {
		cfg.SMTP.Sender = os.Getenv("SMTP_SENDER")
	}
	if os.Getenv("SMTP_RECEIVER") != "" {
		cfg.SMTP.Receiver = os.Getenv("SMTP_RECEIVER")
	}
}

func validateConfig(cfg *config) {
	if cfg.Settings.Sensitivity <= 0.0 || cfg.Settings.Sensitivity > 1.0 {
		log.Fatal("Sensitivity is out of range 0.0 - 1.0\n")
	}
	if cfg.Settings.KeepThreshold <= 0.0 || cfg.Settings.KeepThreshold > 1.0 {
		log.Fatal("KeepThreshold is out of range 0.0 - 1.0\n")
	}
	if cfg.Settings.UploadThreshold <= 0.0 || cfg.Settings.UploadThreshold > 1.0 {
		log.Fatal("UploadThreshold is out of range 0.0 - 1.0\n")
	}
	if cfg.Settings.EmailThreshold <= 0.0 || cfg.Settings.EmailThreshold > 1.0 {
		log.Fatal("EmailThreshold is out of range 0.0 - 1.0\n")
	}
	if cfg.Settings.BaseImageDir == "" {
		log.Fatal("BaseImageDir must be defined\n")
	}
	if cfg.Settings.LogFile == "" {
		log.Fatal("LogFile must be defined\n")
	}
	if cfg.Settings.FFmpegCmd == "" {
		log.Fatal("FFmpegCmd must be defined\n")
	}
}
