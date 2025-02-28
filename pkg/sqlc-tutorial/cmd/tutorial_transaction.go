package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	db "tutorial.sqlc.dev/app/db/codegen/schema"
)

func transaction(ctx context.Context, conn *pgx.Conn, queries *db.Queries) error {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		fmt.Printf("Unable to begin transaction: %v\n", err)
	}

	authorName := "Reverting on error"
	_, err = queries.WithTx(tx).CreateAuthor(ctx, db.CreateAuthorParams{
		Name: authorName,
	})
	if err != nil {
		return tx.Rollback(ctx)
	}

	_, err = queries.WithTx(tx).AddBook(ctx, db.AddBookParams{
		Title:    "Wrong user ID, should revert",
		AuthorID: 404,
	})
	if err != nil {
		fmt.Printf("Expecting to be unable to add a book with %v\n", err)
		_ = tx.Rollback(ctx)
	}

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	for _, author := range authors {
		if author.Name == authorName {
			panic(fmt.Sprintf("Author '%s' should not be present", authorName))
		}
	}

	return nil
}
