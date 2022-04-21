build:
	cd cmd/gs/; go build -o ../../bin/gs .

run:
	go run coscanner/pkg

clean:
	rm -rf bin/*
