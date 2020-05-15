#!make
include .env
export $(shell sed 's/=.*//' .env)

test:
	go test -v github.com/homemade/facecloth