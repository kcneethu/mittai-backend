compile-linux:
	echo "Compiling for Linux OS"
	go env -w GOOS=linux
	go env CC=gcc
	go env -w CGO_ENABLED=1
	go build -o linux/mittai
compile-windows:
	echo "Compiling for Windows OS"
	go env -w GOOS=windows
	go env -w CGO_ENABLED=1
	go build -o windows/mittai.exe

compile-mac:
	echo "Compiling for MacOS arm64"
	go env -w GOOS=darwin
	go env -w GOARCH=arm64
	go env -w CGO_ENABLED=1
	go build -o mac/mittai

clean:
	rm -f linux/mittai windows/mittai.exe mac/mittai

all: compile-linux compile-windows compile-mac