SHELL := /bin/bash

TAG := 1.15.0
APP_NAME := fimpgo

ARCH ?= armhf

BIN_DIR := ./build
TARGET_BIN := $(BIN_DIR)/$(APP_NAME)_$(TAG)_$(ARCH)
MAIN_SRC := ./cli/client.go

all: build-arm

clean:
	-rm -f $(BIN_DIR)/*
	-rm -f $(APP_NAME)
	-rm -f $(TARGET_BIN)
	mkdir -p $(BIN_DIR)

build: clean
	go build -ldflags="-s -w -X main.Version=$(TAG)" -o $(APP_NAME)

build-arm: ARCH=armhf
build-arm: clean
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w -X main.Version=$(TAG)" -o $(TARGET_BIN) $(MAIN_SRC)

build-amd64: ARCH=amd64
build-amd64: clean
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(TAG)" -o $(TARGET_BIN) $(MAIN_SRC)

build-mac: ARCH=amd64
build-mac: clean
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(TAG)" -o $(TARGET_BIN) $(MAIN_SRC)

test:
	go test ./... -v

.PHONY: all clean test build build-arm build-amd64 build-mac
