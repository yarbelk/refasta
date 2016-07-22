GIT_VERSION := $(shell git describe --abbrev=4 --dirty --always --tags)

all: refasta.linux refasta.darwin refasta.exe
.PHONY: all

refasta.linux: main.go
	GOOS=linux go build -ldflags "-X main.version=$(GIT_VERSION)" -o $@ $< 

refasta.darwin: main.go
	GOOS=darwin go build -ldflags "-X main.version=$(GIT_VERSION)" -o $@ $<

refasta.exe: main.go
	GOOS=windows go build -ldflags "-X main.version=$(GIT_VERSION)" -o $@ $<

.PHONY: clean
clean:
	rm refasta.linux refasta.darwin refasta.exe
