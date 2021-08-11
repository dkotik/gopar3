default: test
	go run ./cmd/gopar3/ --growth 3 encode --fragments 55 targetOne --growth 2
test:
	# go test -failfast -v -count 1 -run ^ExampleScanner
	for s in $$(go list ./...); do if ! go test -failfast -v -p 1 $$s; then break; fi; done











.PHONY: default test
