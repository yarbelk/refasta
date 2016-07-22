all: refasta.linux refasta.darwin refasta.exe
.PHONY: all

refasta.linux: main.go
	GOOS=linux go build -ldflags "-X main.version=0.2.0" -o $@ $< 

refasta.darwin: main.go
	GOOS=darwin go build -ldflags "-X main.version=0.2.0" -o $@ $<

refasta.exe: main.go
	GOOS=windows go build -ldflags "-X main.version=0.2.0" -o $@ $<

.PHONY: clean
clean:
	rm refasta.linux refasta.darwin refasta.exe
