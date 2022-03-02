build:
	echo "Compiling binaries for Linux, Mac-Inter, Mac-M1 and Windows"
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/ukrainian-warship-mac-M1
    GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/ukrainian-warship-mac
    GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/ukrainian-warship-linux
    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/ukrainian-warship.exe
