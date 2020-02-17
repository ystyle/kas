export GO111MODULE=on
GOOS=windows go build -ldflags "-s -w" -o kas.exe main.go

GOOS=linux go build -ldflags "-s -w" -o kas-linux main.go

GOOS=darwin go build -ldflags "-s -w" -o kas-darwin main.go

rice append --exec kas.exe
rice append --exec kas-linux
rice append --exec kas-darwin

echo "done!"