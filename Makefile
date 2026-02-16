SHELL := /bin/bash

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
APP_NAME := fimpgo

ARCH ?= armhf

BIN_DIR := ./build
TARGET_BIN := $(BIN_DIR)/$(APP_NAME)_$(VERSION)_$(ARCH)
MAIN_SRC := ./cli/client.go

all: build-arm

configure:
	-rm -f $(BIN_DIR)/*
	-rm -f $(APP_NAME)
	-rm -f $(TARGET_BIN)
	mkdir -p $(BIN_DIR)

build: configure
	go build -ldflags="-s -w" -o $(APP_NAME) $(MAIN_SRC)

build-arm: ARCH=armhf
build-arm: configure
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o $(TARGET_BIN) $(MAIN_SRC)

build-amd64: ARCH=amd64
build-amd64: configure
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(TARGET_BIN) $(MAIN_SRC)

build-mac: ARCH=amd64
build-mac: configure
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(TARGET_BIN) $(MAIN_SRC)

test:
	go test ./...

.PHONY: all configure test build build-arm build-amd64 build-mac
