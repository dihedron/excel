package query

import (
	"log/slog"

	"github.com/dihedron/excel/model"
	"github.com/jmoiron/sqlx"
)

type Query struct {
	Positional struct {
		Query string `positional-arg-name:"query" description:"The SQL query to execute." required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

func (cmd *Query) Execute(args []string) error {
	slog.Debug("querying the database")

	// Connect using sqlx (wraps standard database/sql)
	db, err := sqlx.Connect("sqlite3", "excel.db")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	defer db.Close()

	var facts []model.Fatto
	err = db.Select(&facts, cmd.Positional.Query)
	if err != nil {
		slog.Error("failed to query database", "error", err)
		return err
	}
	return nil
}
