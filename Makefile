build: test
	cd cmd/gs/; go build -o ../../bin/gs .

test:
	go test ./...	

run:
	go run coscanner/pkg

clean:
	rm -rf bin/*
