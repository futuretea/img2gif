all: build
.PHONY: init vendor build linux windows test clean
init:
	go mod init img2gif
vendor:
	go mod vendor
build: linux windows
linux:
	GOOS=linux GOARCH=amd64 go build
windows:
	GOOS=windows GOARCH=amd64 go build
test:
	./img2gif -p ./test
clean:
	rm -f img2gif img2gif.exe output.gif
