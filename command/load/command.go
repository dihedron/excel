package load

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/jszwec/csvutil"
	_ "github.com/mattn/go-sqlite3"
)

type Load struct {
	CSV string `short:"c" long:"csv" description:"The CSV file to load." required:"true"`
}

type Record struct {
	ID    int    `csv:"id" db:"id"`
	Email string `csv:"email" db:"email"`
}

const schema = `
CREATE TABLE IF NOT EXISTS record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT
);`

func (cmd *Load) Execute(args []string) error {

	slog.Debug("opening CVS file", "file", cmd.CSV)
	file, err := os.OpenFile(cmd.CSV, os.O_RDONLY, 0)
	if err != nil {
		slog.Error("error opening CSV file", "error", err)
		return err
	}
	defer file.Close()

	slog.Debug("opening database")
	db, err := sqlx.Connect("sqlite3", "records.db")
	if err != nil {
		slog.Error("error opening database", "error", err)
		return err
	}
	defer db.Close()

	slog.Debug("creating schema")
	if _, err = db.Exec(schema); err != nil {
		slog.Error("error creating schema", "error", err)
		return err
	}

	csvReader := csv.NewReader(file)

	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		slog.Error("error decoding CSV", "error", err)
		return err
	}

	var records []Record
	for {
		var r Record
		if err := dec.Decode(&r); err != nil {
			if err.Error() == "EOF" {
				break
			}
			slog.Error("error decoding CSV", "error", err)
			return err
		}
		records = append(records, r)
	}

	fmt.Println(records)

	return nil
}
