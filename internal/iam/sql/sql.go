package sql

import (
	"context"
	"embed"

	"github.com/adoublef/daw/sql"
)

//go:embed all:*.up.sql
var fsys embed.FS

func Up(ctx context.Context, db *sql.Conn) error {
	return sql.Up(ctx, db, fsys)
}