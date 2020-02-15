export GO111MODULE=on
GOOS=windows go build -ldflags "-s -w" -o hcc.exe main.go

GOOS=linux go build -ldflags "-s -w" -o hcc-linux main.go

GOOS=darwin go build -ldflags "-s -w" -o hcc-darwin main.go

echo "done!"