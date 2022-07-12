demo:
	go build -o bin/demo github.com/runaek/clap/examples/demo

mocks:
	go generate arg.go