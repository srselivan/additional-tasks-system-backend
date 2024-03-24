migrate-create:
	migrate create -ext sql -dir ./migrations -seq -digits 6 initial