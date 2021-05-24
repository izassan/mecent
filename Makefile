build:
	go build ./cmd/fsender/main.go

help:
	go run ./cmd/fsender/main.go -h

run:
	go run ./cmd/fsender/main.go ./cmd/fsender/convert.go 
