.PHONY: all
all: bin/tar

go.mod:
	go mod tidy

bin/tar: *.go go.mod
	go build -o bin/tar main.go

.PHONY: test
test:
	go test -shuffle on -cover -count 10 ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: errcheck
errcheck:
	errcheck ./...

.PHONY: staticcheck
staticcheck:
	staticcheck -checks="all,-ST1000" ./...

.PHONY: clean
clean:
	rm -rf bin/*
