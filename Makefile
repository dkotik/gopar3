default:
	go run ./cmd/gopar3/ --growth 3 encode --fragments 55 targetOne --growth 2
test:
	go test -v ./...
