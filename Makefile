codegen:
	go build -o bin/codegen github.com/runaek/clap/cmd/codegen

eg:
	go build -o bin/example_simple github.com/runaek/clap/examples/derived && \
	go build -o bin/example_derived github.com/runaek/clap/examples/explicit && \
	go build -o bin/example_readme github.com/runaek/clap/examples/readme

generate-derivers: codegen
	./bin/codegen int int --parser=parse.Int --output=pkg/derive/int.go && \
	./bin/codegen ints []int --parser=parse.Ints --output=pkg/derive/ints.go && \
	./bin/codegen float float64 --parser=parse.Float64 --output=pkg/derive/float.go && \
	./bin/codegen floats float64 --parser=parse.Float64 --slice --output pkg/derive/floats.go && \
	./bin/codegen bool bool --parser=parse.Bool --output=pkg/derive/bool.go -F && \
	./bin/codegen counter parse.C --parser=parse.Counter --output=pkg/derive/counter.go -F && \
	./bin/codegen indicator parse.I --parser=parse.Indicator --output=pkg/derive/indicator.go -F

generate-mocks:
	go generate arg.go

test:
	go test -v -tags clap_mocks ./...

testc:
	go test -json -tags clap_mocks ./... -covermode=atomic -coverprofile test.cov

lint:
	go fmt ./... && golangci-lint run