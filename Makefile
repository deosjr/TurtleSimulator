mp4:
	make go && ffmpeg -i turtle.avi turtle.mp4

go:
	go run main.go turtle.go visualisation.go
