package cfg

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type ConfigCamera struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type ConfigFTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type ConfigSMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Sender   string `yaml:"sender"`
	Receiver string `yaml:"receiver"`
}

type ConfigSettings struct {
	Sensitivity     float32 `yaml:"sensitivity"`
	KeepThreshold   float32 `yaml:"keepThreshold"`
	UploadThreshold float32 `yaml:"uploadThreshold"`
	EmailThreshold  float32 `yaml:"emailThreshold"`
	EmailInterval   int     `yaml:"emailInterval"`
	BaseImageDir    string  `yaml:"baseImageDir"`
	LogFile         string  `yaml:"logFile"`
	FFmpegCmd       string  `yaml:"ffmpegCmd"`
}

// Config contains configuration
type Config struct {
	Camera   ConfigCamera   `yaml:"camera"`
	FTP      ConfigFTP      `yaml:"ftp"`
	SMTP     ConfigSMTP     `yaml:"smtp"`
	Settings ConfigSettings `yaml:"settings"`
}

var (
	// Camera configuration
	Camera ConfigCamera
	// FTP configuration
	FTP ConfigFTP
	// SMTP configuration
	SMTP ConfigSMTP
	// Settings configuration
	Settings ConfigSettings
)

// LoadConfig loads configuration from yml file or default
func LoadConfig(configFile string) {
	cfg := Config{}
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

	loadEnvSecrets(&cfg)
	validateConfig(&cfg)

	Camera = cfg.Camera
	FTP = cfg.FTP
	SMTP = cfg.SMTP
	Settings = cfg.Settings
}

func loadEnvSecrets(cfg *Config) {
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

func validateConfig(cfg *Config) {
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
	if cfg.Settings.EmailInterval < 0 || cfg.Settings.EmailInterval > 3600 {
		log.Fatal("EmailInterval is out of range 0 - 3600 seconds\n")
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
