package main

import (
	"context"
	"log"
	"reflect"
	db "tutorial.sqlc.dev/app/db/codegen/schema"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func run(dsn string) error {
	ctx := context.Background()
	if dsn == "" {
		dsn = "postgres://pqgotest:pqgotest@localhost:5432/pqgotest"
	}

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  pgtype.Text{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}
