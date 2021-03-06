# Watchdog
Camera alert system which can run on Raspberry Pi, capture still images from RTSP enabled camera stream and send alerts to FTP and e-mail via SMTP.

## Features
- Camera connectivity using ffmpeg and rtsp protocol capturing still images
- Upload to FTP triggered by threshold
- Email triggered by threshold
- Rotating logs

## Prerequisites
- Camera supporting RTSP protocol
- ffmpeg
- ftp account
- smtp email account

## Configuration
Update config.yml, as a best practice do not put your secret credentials into this file.
Create a .secrets file with credentials:
```sh
WATCHDOG_ID=
LOG_FILE=
IMAGE_DIR=
CAMERA_HOST=
CAMERA_PORT=554
CAMERA_USER=
CAMERA_PASS=

FTP_HOST=
FTP_PORT=990
FTP_USER=
FTP_PASS=

SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASS=
SMTP_SENDER=
SMTP_RECEIVER=
```

## Build
```sh
make build
make pi-build
```

## Run
First edit Makefile, config.yml and .secrets to ensure you have proper settings for your environment.
```sh
make run
```

## Docker
First edit Makefile, config.yml and .secrets to ensure you have proper settings for your environment.
Also ensure that DOCKER_IMAGE_DIR and DOCKER_LOG_DIR point to existing absolute path.
```sh
make docker-build
make docker-run
```