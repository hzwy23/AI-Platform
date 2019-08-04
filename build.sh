#!/bin/bash
echo $1
rmdir "application"
mkdir "application"

if [ "$1" = "windows" ];then 
	echo "编译windows平台软件"
	GOOS=windows GOARCH=amd64 go build -o ai_lamp.exe main.go
fi 

if [ "$1" = "macos" ];then 
	echo "编译$1平台软件"
	GOOS=macos GOARCH=amd64 go build -o ai_lamp main.go
fi 

if [ "$1" = "linux" ];then 
	echo "编译$1平台软件"
	GOOS=linux GOARCH=amd64 go build -o ai_lamp main.go
fi 



