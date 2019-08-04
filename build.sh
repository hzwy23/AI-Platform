#!/bin/bash

function clean_application() {
  echo "清除历史安装包"
	rm -rf ./application
	mkdir ./application
}

function pack_application(){
  echo "打包程序到安装包$1"
	mv ./$1 ./application/
}

function copy_static() {
  cp -r ./conf ./application/
	cp -r ./webui ./application/
}

function window() {
  echo "编译windows平台软件"
	GOOS=windows GOARCH=amd64 go build -o windows_start.exe main.go
	pack_application "windows_start.exe"
}

function macos() {
  GOOS=darwin GOARCH=amd64 go build -o macos_start main.go
	pack_application "macos_start"
}

function linux() {
  echo "编译$1平台软件"
	GOOS=linux GOARCH=amd64 go build -o linux_start main.go
		pack_application "linux_start"

}

case $1 in
  "window")
    clean_application
    window
    copy_static
    ;;
  "macos")
    clean_application
    macos
    copy_static
    ;;
  "linux")
    clean_application
    linux
    copy_static
    ;;
  *)
    clean_application
    window
    macos
    linux
    copy_static
    ;;
esac
