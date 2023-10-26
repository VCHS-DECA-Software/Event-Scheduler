codegen-db:
	rm -f schema.db
	rm -rf lib/db
	go run cmd/generate-db/main.go --db=schema.db
	jet -dsn=file://$(shell pwd)/schema.db -path=./lib/db
