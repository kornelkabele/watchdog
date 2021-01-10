# Watchdog
Camera alert system which successfuly runs on Raspberry Pi

## Features
- Camera connectivity using ffmpeg and rtsp protocol capturing still images
- Upload to FTP triggered by threshold
- Email triggered by threshold

## Configuration

Create a .secrets file with credentials:
```sh
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
`sh
make build
`
`sh
make build-pi
`

## Run

`sh
make run
`
