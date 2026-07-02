package query

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dihedron/excel/command/base"
	"github.com/dihedron/excel/encoder"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-isatty"
)

type Query struct {
	base.Command
	Positional struct {
		Statement string `positional-arg-name:"statement" description:"The SQL query to execute." required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

func (cmd *Query) Execute(args []string) error {
	slog.Debug("querying the database")

	// sonnect using sqlx (wraps standard database/sql)
	db, err := sqlx.Connect("sqlite3", cmd.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	defer db.Close()

	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(cmd.Positional.Statement)), "select") {

		slog.Debug("read-only query", "statement", cmd.Positional.Statement)

		rows, err := db.Query(cmd.Positional.Statement)
		if err != nil {
			slog.Error("failed to query database", "query", cmd.Positional.Statement, "error", err)
			return err
		}
		defer rows.Close()

		// 1. get the column names from the query result
		columns, err := rows.Columns()
		if err != nil {
			slog.Error("failed to retrieve table metadata", "error", err)
			return err
		}

		var entities []map[string]any

		for rows.Next() {
			// 2. create a slice of any's to represent each column
			values := make([]any, len(columns))

			// 3. create a second slice to contain pointers to each item in the values slice
			pointers := make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}

			// 4. scan the result into the column pointers using variadic expansion
			if err := rows.Scan(pointers...); err != nil {
				slog.Error("error scanning values into pointers", "error", err)
				return err
			}

			// 5. build the map for this specific row
			entity := make(map[string]any)
			for i, column := range columns {
				// dereference the pointer to get the actual value
				val := pointers[i].(*any)

				// handle NULL values from the database
				if *val == nil {
					entity[column] = nil
					continue
				}

				// type assertion: many drivers return strings/varchars as []byte;;
				// we convert them to standard Go strings for easier use
				if b, ok := (*val).([]byte); ok {
					entity[column] = string(b)
				} else {
					entity[column] = *val
				}
			}

			slog.Debug("entity ready", "entity", entity)

			entities = append(entities, entity)
		}

		// check for errors encountered during iteration
		if err := rows.Err(); err != nil {
			slog.Error("there were errors retrieving results", "error", err)
			return err
		}

		if isatty.IsTerminal(os.Stdout.Fd()) {
			green := color.New(color.FgGreen).SprintfFunc()
			if len(entities) == 1 {
				fmt.Fprintf(os.Stdout, "OK: %s entity retrieved\n", green(fmt.Sprintf("%d", len(entities))))
			} else {
				fmt.Fprintf(os.Stdout, "OK: %s entities retrieved\n", green(fmt.Sprintf("%d", len(entities))))
			}
		} else {
			if len(entities) == 1 {
				fmt.Fprintf(os.Stdout, "OK: %d entity retrieved\n", len(entities))
			} else {
				fmt.Fprintf(os.Stdout, "OK: %d entities retrieved\n", len(entities))
			}
		}

		// entities ready
		e, err := encoder.New(cmd.Format, encoder.WithIndentation() /* TODO: implement mapper */)
		if err != nil {
			slog.Error("failed to create encoder", "format", cmd.Format, "error", err)
			return err
		}
		defer e.Close()

		for _, entity := range entities {
			slog.Debug("encoding entity to output", "entity", entity, "format", cmd.Format)
			if err := e.Encode(os.Stdout, entity); err != nil {
				slog.Error("failed to encode entity", "error", err)
				return err
			}
		}

	} else {
		var result sql.Result
		result, err = db.Exec(cmd.Positional.Statement)
		if err != nil {
			slog.Error("failed to perform statement", "statement", cmd.Positional.Statement, "error", err)
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if isatty.IsTerminal(os.Stdout.Fd()) {
			green := color.New(color.FgGreen).SprintfFunc()
			fmt.Printf("OK: %s rows affected\n", green(fmt.Sprintf("%d", affected)))
		} else {
			fmt.Printf("OK: %d rows affected\n", affected)
		}
	}
	return nil
}

/*
// QueryDynamic executes a query and returns the entities as a slice of maps,
// where the map keys are column names and the values are the row data.
func QueryDynamic2(db *sql.DB, query string) ([]map[string]any, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 1. Get the column names from the query result
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var entities []map[string]any

	for rows.Next() {
		// 2. Create a slice of any's to represent each column
		columns := make([]any, len(columns))

		// 3. Create a second slice to contain pointers to each item in the columns slice
		pointers := make([]any, len(columns))
		for i := range columns {
			pointers[i] = &columns[i]
		}

		// 4. Scan the result into the column pointers using variadic expansion
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		// 5. Build the map for this specific row
		entity := make(map[string]any)
		for i, column := range columns {
			// Dereference the pointer to get the actual value
			val := pointers[i].(*any)

			// Handle NULL values from the database
			if *val == nil {
				entity[column] = nil
				continue
			}

			// Type assertion: Many drivers return strings/varchars as []byte.
			// We convert them to standard Go strings for easier use.
			if b, ok := (*val).([]byte); ok {
				entity[column] = string(b)
			} else {
				entity[column] = *val
			}
		}

		entities = append(entities, entity)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
*/
