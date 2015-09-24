migrate:
	migrate -url $DATABASE_URI?sslmode=disable -path ./migrations/ up
