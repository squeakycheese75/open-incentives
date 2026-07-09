up:
	docker compose up --build

reset:
	rm -f data/app.db
	docker compose up --build

clean:
	rm -f data/app.db
	DATABASE_SEED=false docker compose up --build

test:
	go test -short ./...

lint:
	golangci-lint version && golangci-lint run --verbose  -E  misspell    