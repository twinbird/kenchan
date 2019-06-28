PROGRAM=kenchan

kenchan:
	go build

.PHONY: clean
clean:
	rm -f kenchan
	rm -rf bin
	rm -rf KEN_ALL.CSV

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux64/${PROGRAM}
	GOOS=linux GOARCH=386 go build -o ./bin/linux386/${PROGRAM}
	GOOS=windows GOARCH=386 go build -o ./bin/windows386/${PROGRAM}.exe
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows64/${PROGRAM}.exe
	GOOS=darwin GOARCH=386 go build -o ./bin/darwin386/${PROGRAM}
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin64/${PROGRAM}
