package postgres

import (
	"fmt"
	"context"

	"github.com/go-pg/pg/v10"
)

type DBLogger struct {}

func (d DBLogger) BeforeQuery(ctx context.Context ,q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d DBLogger) AfterQuery(ctx context.Context ,q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	fmt.Println(string(fq))
	return nil
}

// New creates a new DB instance
func New(opts *pg.Options) *pg.DB {
	db := pg.Connect(opts)

	return db
}