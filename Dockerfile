# build image
FROM golang:alpine AS builder
LABEL stage=builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o watchdog -v ./cmd/watchdog/main.go


# final stage
FROM alpine AS final
RUN apk add --no-cache ffmpeg
WORKDIR /
COPY --from=builder /src/watchdog .
COPY --from=builder /src/config.yml .

# executable
CMD [ "./watchdog" ]