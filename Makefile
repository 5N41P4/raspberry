.PHONY: app/test ui/build fetch

raspberry:
	go build -o=./bin/raspberry ./cmd/web

app/test: ui/build
	go run ./cmd/web

ui/build: 
	sudo cp -r ~/Source/pineapple/dist /usr/local/raspberry/ui

fetch:
	git pull