migrate:
	migrate create -seq -ext .sql -dir ./db/migrations ${name}

up:
	 migrate -path ./db/migrations -database ${DB_DSN} up

down:
	 migrate -path ./db/migrations -database ${DB_DSN} down ${q}