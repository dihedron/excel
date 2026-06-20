package query

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/dihedron/excel/encoder"
	"github.com/dihedron/excel/model"
	"github.com/jmoiron/sqlx"
)

type Query struct {
	Format     string `short:"t" long:"format" description:"The format of the output." optional:"true" default:"text" choice:"text" choice:"json" choice:"yaml" choice:"csv"`
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

	e, err := encoder.New(cmd.Format, encoder.WithIndentation(), encoder.WithDataMapper(func(data any) ([]string, error) {
		f := data.(model.Fatto)
		return []string{f.Anno, f.CID, f.CodiceIndividuale, f.Nome, f.Cognome, f.Dipartimento, f.Servizio, f.Divisione, f.Settore, f.Segmento.String(), f.DecorrenzaSegmento.Format("02/01/2006"), strconv.Itoa(f.Livello), f.DecorrenzaLivello.Format("02/01/2006")}, nil
	}))
	if err != nil {
		slog.Error("failed to create encoder", "format", cmd.Format, "error", err)
		return err
	}
	defer e.Close()

	for _, fact := range facts {

		if err := e.Encode(os.Stdout, fact); err != nil {
			slog.Error("failed to encode fatto", "error", err)
			return err
		}
	}

	return nil
}
