package query

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dihedron/excel/encoder"
	"github.com/jmoiron/sqlx"
)

type Query struct {
	Format     string `short:"t" long:"format" description:"The format of the output." optional:"true" default:"text" choice:"text" choice:"json" choice:"yaml" choice:"csv" env:"EXCEL_FORMAT"`
	Positional struct {
		Statement string `positional-arg-name:"statement" description:"The SQL query to execute." required:"yes"`
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

	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(cmd.Positional.Statement)), "select") {

		rows, err := db.Query(cmd.Positional.Statement)
		if err != nil {
			return err
		}
		defer rows.Close()

		// 1. Get the column names from the query result
		cols, err := rows.Columns()
		if err != nil {
			return err
		}

		var results []map[string]any

		for rows.Next() {
			// 2. Create a slice of any's to represent each column
			columns := make([]any, len(cols))

			// 3. Create a second slice to contain pointers to each item in the columns slice
			columnPointers := make([]any, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			// 4. Scan the result into the column pointers using variadic expansion
			if err := rows.Scan(columnPointers...); err != nil {
				return err
			}

			// 5. Build the map for this specific row
			rowMap := make(map[string]any)
			for i, colName := range cols {
				// Dereference the pointer to get the actual value
				val := columnPointers[i].(*any)

				// Handle NULL values from the database
				if *val == nil {
					rowMap[colName] = nil
					continue
				}

				// Type assertion: Many drivers return strings/varchars as []byte.
				// We convert them to standard Go strings for easier use.
				if b, ok := (*val).([]byte); ok {
					rowMap[colName] = string(b)
				} else {
					rowMap[colName] = *val
				}
			}

			results = append(results, rowMap)
		}

		// Check for errors encountered during iteration
		if err := rows.Err(); err != nil {
			return err
		}

		// results ready

		e, err := encoder.New(cmd.Format, encoder.WithIndentation() /* TODO: implement mapper */)
		if err != nil {
			slog.Error("failed to create encoder", "format", cmd.Format, "error", err)
			return err
		}
		defer e.Close()

		for _, result := range results {
			if err := e.Encode(os.Stdout, result); err != nil {
				slog.Error("failed to encode fatto", "error", err)
				return err
			}
		}

	} else {
		var result sql.Result
		result, err = db.Exec(cmd.Positional.Statement)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		fmt.Printf("statement OK, %d rows affected\n", affected)

		rows, err := db.Query("SELECT count(*) FROM fatti where dipartimento = 'DIPARTIMENTO INDFORMATICA'")
		if err != nil {
			return err
		}
		defer rows.Close()
	}

	// var facts []model.Fatto
	//  err = db.Select(&facts, cmd.Positional.Query)
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

// QueryDynamic executes a query and returns the results as a slice of maps,
// where the map keys are column names and the values are the row data.
func QueryDynamic(db *sql.DB, query string) ([]map[string]any, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 1. Get the column names from the query result
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]any

	for rows.Next() {
		// 2. Create a slice of any's to represent each column
		columns := make([]any, len(cols))

		// 3. Create a second slice to contain pointers to each item in the columns slice
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// 4. Scan the result into the column pointers using variadic expansion
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// 5. Build the map for this specific row
		rowMap := make(map[string]any)
		for i, colName := range cols {
			// Dereference the pointer to get the actual value
			val := columnPointers[i].(*any)

			// Handle NULL values from the database
			if *val == nil {
				rowMap[colName] = nil
				continue
			}

			// Type assertion: Many drivers return strings/varchars as []byte.
			// We convert them to standard Go strings for easier use.
			if b, ok := (*val).([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = *val
			}
		}

		results = append(results, rowMap)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
