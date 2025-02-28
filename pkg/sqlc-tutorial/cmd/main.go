package main

import (
	"log"
	"tutorial.sqlc.dev/app/db/migrations"
)

func main() {
	dsn := "postgres://pqgotest:pqgotest@localhost:5432/pqgotest?sslmode=disable"
	migrations.Run(dsn)
	if err := run(dsn); err != nil {
		log.Fatal(err)
	}
}
