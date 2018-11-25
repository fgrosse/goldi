.PHONY: all test install test coverage
all: test

install:
	go get golang.org/x/tools/cmd/cover
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls
	go get -v -t ./...

test:
	go test -cover ./...

coverage:
	go test -cover -coverpkg github.com/fgrosse/goldi -coverprofile goldi.coverprofile .
	go test -cover -coverpkg github.com/fgrosse/goldi/goldigen -coverprofile goldigen.coverprofile ./goldigen
	go test -cover -coverpkg github.com/fgrosse/goldi/validation -coverprofile validation.coverprofile ./validation
