
gen:
	go build -o bin/codegen github.com/runaek/clap/cmd/generate_derivers

generate: gen
	./bin/codegen int int --parser=parse.Int --output=pkg/derive/int.go && \
	./bin/codegen ints []int --parser=parse.Ints --output=pkg/derive/ints.go && \
	./bin/codegen float float64 --parser=parse.Float64 --output=pkg/derive/float.go && \
	./bin/codegen counter parse.C --parser=parse.Counter --output=pkg/derive/counter.go

test:
	go test -v -tags clap_mocks ./...

testc:
	go test -json -tags clap_mocks ./... -covermode=atomic -coverprofile test.cov

mocks:
	go generate arg.go