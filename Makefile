all: install

GOPATH:=$(CURDIR)
export GOPATH

dep:
	go get github.com/go-sql-driver/mysql

install:dep
	go install dbstatus
