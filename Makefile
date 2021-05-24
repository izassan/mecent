test_fsender:
	go run ./cmd/fsender/main.go ./cmd/fsender/convert.go

install_fsender:
	go install ./cmd/fsender
