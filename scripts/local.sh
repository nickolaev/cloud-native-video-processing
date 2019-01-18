#!/bin/sh

function cleanup {
    kill -9 ${PID_SOURCE} ${PID_LOGO} ${PID_TIMESTAMP} >/dev/null 2>&1
}
trap cleanup EXIT

go run cmd/streamer/*.go -source ./assets/video.mp4 -port 2000 &
PID_SOURCE=$!

go run cmd/streamer/*.go -source http://localhost:2000 -logo ./assets/logo.png -port 2001 &
PID_LOGO=$!

go run cmd/streamer/*.go -source http://localhost:2001 -font ./assets/bedstead.otf -timestamp -port 2002 &
PID_TIMESTAMP=$!

go run cmd/streamer/*.go -source http://localhost:2002 -preset HD -port 2003
