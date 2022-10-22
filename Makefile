mp4:
	make go && ffmpeg -i turtle.avi turtle.mp4

run:
	go run main.go visualisation.go

gen:
	go generate ./...

test:
	go test ./...
