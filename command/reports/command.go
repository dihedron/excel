package reports

import (
	"log/slog"

	"github.com/dihedron/excel/logic"
	"github.com/jmoiron/sqlx"
)

type Reports struct {
	// Format     string `short:"t" long:"format" description:"The format of the output." optional:"true" default:"text" choice:"text" choice:"json" choice:"yaml" choice:"csv"`
	// Positional struct {
	// 	Query string `positional-arg-name:"query" description:"The SQL query to execute." required:"yes"`
	// } `positional-args:"yes" required:"yes"`
}

func (cmd *Reports) Execute(args []string) error {
	slog.Debug("generating reports")

	// Connect using sqlx (wraps standard database/sql)
	db, err := sqlx.Connect("sqlite3", "excel.db")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	defer db.Close()

	err = logic.Avanzamenti(db)
	if err != nil {
		slog.Error("failed to generate reports", "error", err)
		return err
	}

	// var facts []model.Fatto
	// err = db.Select(&facts, cmd.Positional.Query)
	// if err != nil {
	// 	slog.Error("failed to query database", "error", err)
	// 	return err
	// }

	// e, err := encoder.New(cmd.Format, encoder.WithIndentation(), encoder.WithDataMapper(model.MapFattoToSlice))
	// if err != nil {
	// 	slog.Error("failed to create encoder", "format", cmd.Format, "error", err)
	// 	return err
	// }
	// defer e.Close()

	// for _, fact := range facts {
	// 	if err := e.Encode(os.Stdout, fact); err != nil {
	// 		slog.Error("failed to encode fatto", "error", err)
	// 		return err
	// 	}
	// }

	return nil
}
