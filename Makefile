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

generate-mocks:
	@echo "Generating mocks..."
	mockgen -source=internal/eval/usecase/evaluate/main.go -destination=internal/eval/usecase/evaluate/mocks/mocks.go -package=mocks

admin-dev:
	cd apps/admin && npm install && npm run dev

admin-lint:
	cd apps/admin && npm run lint
