codegen-db:
	go run cmd/generate-db/main.go --db=schema.db
	jet -dsn=file://$(shell pwd)/schema.db -path=./lib/db
	rm schema.db
