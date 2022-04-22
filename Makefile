build:
	cd cmd/gs/; go build -o ../../bin/gs .

test:
	go test ./...	

clean:
	rm -rf bin/*
