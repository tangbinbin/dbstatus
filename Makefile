all: install

install:
	go build -o bin/dbstatus ./src/dbstatus
