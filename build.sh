#!/bin/bash

if [ -f "main" ]; then
    rm main;
fi

if [ -f "main.zip" ]; then
    rm main.zip
fi

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
zip main.zip main