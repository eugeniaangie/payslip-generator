build:
	go build -o server ./cmd/*go

run:
	go run ./cmd/*go start

watch:
	reflex -s -r '\.go$$' make run
