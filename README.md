# Watchdog
Camera alert system which can run on Raspberry Pi, capture still images from RTSP enabled camera stream and send alerts to FTP and e-mail via SMTP.

## Features
- Camera connectivity using ffmpeg and rtsp protocol capturing still images
- Upload to FTP triggered by threshold
- Email triggered by threshold

## Configuration

Create a .secrets file with credentials:
```sh
WATCHDOG_ID=""
CAMERA_HOST=""
CAMERA_PORT=554
CAMERA_USER=""
CAMERA_PASS=""

FTP_HOST=""
FTP_PORT=990
FTP_USER=""
FTP_PASS=""

SMTP_HOST=""
SMTP_PORT=587
SMTP_USER=""
SMTP_PASS=""
SMTP_SENDER=""
SMTP_RECEIVER=""
```

## Build
```sh
make build
make build-pi
```

## Run

`make run`
