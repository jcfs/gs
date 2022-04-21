build:
	go build -o bin/gs pkg/*.go

run:
	go run coscanner/pkg

clean:
	rm -rf bin/*
