package load

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dihedron/excel/command/base"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

type Load struct {
	base.Command
	File    string    `short:"f" long:"file" description:"The Excel file to load." required:"true"`
	Sheet   string    `short:"s" long:"sheet" description:"The sheet to load." required:"true"`
	Columns []Mapping `short:"m" long:"mapping" description:"The columns mappings (format: <field>:<offset|value>[:converter])." required:"true"`
	//Filters    []Filter  `short:"x" long:"filter" description:"Only output records matching the filter (format: <field>:<value>)." optional:"true"`
	Headers    int `short:"h" long:"headers" description:"The number of headers to skip." optional:"true" default:"0"`
	Positional struct {
		Statement string `positional-arg-name:"statement" description:"The SQL statement to execute." required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

type Mapping struct {
	Name      string
	Offset    int
	Value     string
	Converter func(string) (any, error)
}

func (m *Mapping) UnmarshalFlag(value string) error {
	values := strings.Split(value, ":")
	m.Name = values[0]
	if offset, err := getOffset(values[1]); err == nil && offset > -1 {
		m.Offset = offset
		if len(values) > 2 {
			if values[2] == "int" {
				m.Converter = ToInt()
			} else if strings.HasPrefix(values[2], "time") {
				m.Converter = ToTime(getFormat(values[2]))
			} else if values[2] == "segmento" {
				m.Converter = ToSegmento()
			}
		}
	} else {
		m.Offset = -1
		m.Value = values[1]
	}
	return nil
}

type Filter struct {
	Field string
	Value string
}

func (f *Filter) UnmarshalFlag(value string) error {
	values := strings.Split(value, ":")
	f.Field = values[0]
	f.Value = values[1]
	return nil
}

type Sheet struct {
	Name  string
	Table string
}

func (s *Sheet) UnmarshalFlag(value string) error {
	values := strings.Split(value, ":")
	s.Name = values[0]
	s.Table = values[1]
	return nil
}

func (cmd *Load) Execute(args []string) error {

	//slog.Debug("opening CVS file", "file", cmd.File, "sheet", cmd.Sheet, "mappings", cmd.Columns, "filters", cmd.Filters)
	slog.Debug("opening CVS file", "file", cmd.File, "sheet", cmd.Sheet, "mappings", cmd.Columns)

	// open the Excel file
	f, err := excelize.OpenFile(cmd.File)
	if err != nil {
		slog.Error("failed to open file", "path", cmd.File, "error", err)
		return err
	}
	defer f.Close()

	// connect using sqlx (wraps standard database/sql)
	db, err := sqlx.Connect("sqlite3", cmd.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	defer db.Close()

	ctx := context.Background()

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to open transaction", "error", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, cmd.Positional.Statement)
	if err != nil {
		slog.Error("failed to prepare statement", "statement", cmd.Positional.Statement, "error", err)
		return err
	}
	defer stmt.Close()

	rows, err := f.GetRows(cmd.Sheet)
	if err != nil {
		slog.Error("failed to get rows", "sheet", cmd.Sheet, "error", err)
		return err
	}

	count := 0

rows:
	for i, row := range rows {
		if i < cmd.Headers {
			slog.Debug("skipping line", "count", i, "headers", cmd.Headers)
			continue
		}
		slog.Debug("processing row...", "count", i)
		entity := map[string]any{}
		for _, column := range cmd.Columns {
			if column.Offset == NoOffset {
				slog.Debug("mapping does not represent a column offset, using direct value", "name", column.Name, "value", column.Value)
				if column.Converter != nil {
					slog.Debug("need to convert value first", "name", column.Name, "value", column.Value)
					value, err := column.Converter(column.Value)
					if err != nil {
						slog.Error("failed to convert value", "name", column.Name, "value", column.Value, "error", err)
						fmt.Fprintf(os.Stderr, "%s!%05d: failed to convert value for column %s (value: %q, error: %v)\n", cmd.Sheet, i, column.Name, column.Value, err)
						continue rows
					}
					slog.Debug("value successfully converted", "name", column.Name, "value", value)
					entity[column.Name] = value
				} else {
					slog.Debug("using value as provided", "name", column.Name, "value", column.Value)
					entity[column.Name] = column.Value
				}
			} else {
				// this is an offset into the eXcel file
				slog.Debug("mapping represents a column offset into the eXcel file", "name", column.Name, "offset", column.Offset)
				if column.Converter != nil {
					slog.Debug("need to convert value first", "name", column.Name, "value", column.Value)
					value, err := column.Converter(row[column.Offset])
					if err != nil {
						slog.Error("failed to convert value", "name", column.Name, "value", row[column.Offset], "error", err)
						fmt.Fprintf(os.Stderr, "%s!%05d: failed to convert value for column %s (value: %q, error: %v)\n", cmd.Sheet, i, column.Name, row[column.Offset], err)
						continue rows
					}
					slog.Debug("value successfully converted", "name", column.Name, "value", value)
					entity[column.Name] = value
				} else {
					slog.Debug("using value as provided", "name", column.Name, "value", column.Value)
					entity[column.Name] = row[column.Offset]
				}
			}
		}
		slog.Debug("inserting entity into database", "statement", cmd.Positional.Statement, "entity", entity)
		if _, err := stmt.ExecContext(ctx, entity); err != nil {
			slog.Error("failed to insert entity", "entity", entity, "error", err)
			fmt.Fprintf(os.Stderr, "%s!%05d: failed to insert row %+v into database: %v\n", cmd.Sheet, i, entity, err)
			continue rows
		}
		count++

		// 		match := false
		// 		if len(cmd.Filters) == 0 {
		// 			match = true
		// 		} else {
		// 		filters:
		// 			for _, filter := range cmd.Filters {
		// 				switch filter.Field {
		// 				case "CID":
		// 					if fatto.CID == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Codice Individuale":
		// 					if fatto.CodiceIndividuale == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Nominativo":
		// 					if fatto.Nominativo == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Dipartimento":
		// 					if fatto.Dipartimento == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Servizio":
		// 					if fatto.Servizio == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Divisione":
		// 					if fatto.Divisione == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Categoria":
		// 					if fatto.Settore == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Segmento":
		// 					if fatto.Segmento.String() == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Decorrenza Segmento":
		// 					if fatto.DecorrenzaLivello.Format("02/01/2006") == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Livello":
		// 					value, err := strconv.Atoi(filter.Value)
		// 					if err != nil {
		// 						slog.Error("failed to parse integer", "integer", filter.Value, "error", err)
		// 						return err
		// 					}
		// 					if fatto.Livello == value {
		// 						match = true
		// 						break filters
		// 					}
		// 				case "Decorrenza Livello":
		// 					if fatto.DecorrenzaLivello.Format("02/01/2006") == filter.Value {
		// 						match = true
		// 						break filters
		// 					}
		// 				}
		// 			}
		// 		}

		// e, err := encoder.New(cmd.Format, encoder.WithIndentation(), encoder.WithDataMapper(model.MapFattoToSlice))
		// if err != nil {
		// 	slog.Error("failed to create encoder", "format", cmd.Format, "error", err)
		// 	return err
		// }
		// defer e.Close()

		// if match {
		// 	if err := e.Encode(os.Stdout, fatto); err != nil {
		// 		slog.Error("failed to encode fatto", "error", err)
		// 		return err
		// 	}
		// switch cmd.Format {
		// case "text":
		// 	fmt.Printf("%+v\n", fatto)
		// case "json":
		// 	b, err := json.MarshalIndent(fatto, "", "  ")
		// 	if err != nil {
		// 		return err
		// 	}
		// 	fmt.Println(string(b))
		// case "yaml":
		// 	b, err := yaml.Marshal(fatto)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	fmt.Println(string(b))
		// case "csv":
		// 	if w == nil {
		// 		w = csv.NewWriter(os.Stdout)
		// 		defer w.Flush()
		// 	}
		// 	w.Write([]string{fatto.Anno, fatto.CID, fatto.CodiceIndividuale, fatto.Nome, fatto.Cognome, fatto.Dipartimento, fatto.Servizio, fatto.Divisione, fatto.Settore, fatto.Segmento.String(), fatto.DecorrenzaSegmento.Format("02/01/2006"), strconv.Itoa(fatto.Livello), fatto.DecorrenzaLivello.Format("02/01/2006")})
		// 	w.Flush()
		// }
		// }

		// if _, err := db.NamedExec(cmd.Positional.Statement, entity); err != nil {
		// 	slog.Error("failed to insert fatto", "fatto", entity, "error", err)
		// 	return err
		// }
	}
	fmt.Fprintf(os.Stdout, "%d lines inserted out of %d\n", count, len(rows)-cmd.Headers)

	/*


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

	*/
	return nil
}

func safeParseDate(value string) (time.Time, error) {
	formats := []string{
		"02/01/2006",
		"02-01-06",
		"01-02-06",
	}
	var result error

	for _, f := range formats {
		t, err := time.Parse(f, value)
		if err == nil {
			return t, nil
		}
		result = errors.Join(result, err)
	}
	return time.Time{}, result
}

const NoOffset = -1

func getOffset(value string) (int, error) {
	match := regexp.MustCompile(`\{(\d+)\}`).FindStringSubmatch(value)
	if len(match) > 1 {
		offset, err := strconv.Atoi(match[1])
		if err == nil {
			slog.Debug("mapping represents column offset", "value", value, "offset", offset)
			return offset, nil
		}
		slog.Error("invalid format in mapping", "value", value, "error", err)
		return NoOffset, err
	}
	slog.Debug("mapping does not represnet a column offset")
	return NoOffset, nil
}

func getFormat(value string) string {
	match := regexp.MustCompile(`time\(([0-9\/\.\-:]+)\)`).FindStringSubmatch(value)
	if len(match) > 1 {
		slog.Debug("valid date forma found", "value", value, "format", match[1])
		return match[1]
	}
	slog.Debug("returning default date format", "value", value, "format", "02/01/2006")
	return "02/01/2006"
}
