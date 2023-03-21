.PHONY: test
test:
	go test -v

.PHONY: coverage
coverage:
	go test -cover 
	

.PHONY: covhtml 
covhtml:
	go test -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html &&	rm -rf c.out
