build:
	go build -o bin/weatherapp

run: build
	./bin/weatherapp