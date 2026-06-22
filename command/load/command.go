package load

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dihedron/excel/encoder"
	"github.com/dihedron/excel/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

type Load struct {
	File    string    `short:"f" long:"file" description:"The Excel file to load." required:"true"`
	Sheets  []Sheet   `short:"s" long:"sheet" description:"The sheet to load (format: <sheet>:<label>)." required:"true"`
	Format  string    `short:"t" long:"format" description:"The format of the output." optional:"true" default:"none" choice:"text" choice:"json" choice:"yaml" choice:"csv" choice:"none"`
	Columns []Mapping `short:"m" long:"mapping" description:"The columns mappings (format: <field>:<offset>)." required:"true"`
	Filters []Filter  `short:"x" long:"filter" description:"Only output records matching the filter (format: <field>:<value>)." optional:"true"`
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
	Label string
}

func (s *Sheet) UnmarshalFlag(value string) error {
	values := strings.Split(value, ":")
	s.Name = values[0]
	s.Label = values[1]
	return nil
}

type Mapping struct {
	Name   string
	Offset int
}

func (m *Mapping) UnmarshalFlag(value string) error {
	values := strings.Split(value, ":")
	m.Name = values[0]
	val, err := strconv.Atoi(values[1])
	if err != nil {
		return fmt.Errorf("error parsing column offset %s", values[1])
	}
	m.Offset = val
	return nil
}

func (cmd *Load) Execute(args []string) error {

	slog.Debug("opening CVS file", "file", cmd.File, "sheets", cmd.Sheets, "mappings", cmd.Columns, "filters", cmd.Filters)

	// open the Excel file
	f, err := excelize.OpenFile(cmd.File)
	if err != nil {
		slog.Error("failed to open file", "path", cmd.File, "error", err)
		return err
	}
	defer f.Close()

	// connect using sqlx (wraps standard database/sql)
	db, err := sqlx.Connect("sqlite3", "excel.db")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	defer db.Close()

	// execute raw schema
	db.MustExec(model.SchemaFatti)

	for _, sheet := range cmd.Sheets {

		rows, err := f.GetRows(sheet.Name)
		if err != nil {
			slog.Error("failed to get rows", "sheet", sheet.Name, "error", err)
			return err
		}

		livello := regexp.MustCompile(`(\d+)`)
	outer:
		for i, row := range rows {
			if i == 0 {
				continue
			}
			anno, err := strconv.Atoi(sheet.Label)
			if err != nil {
				slog.Error("failed to parse integer", "integer", sheet.Label, "error", err)
				return err
			}
			fatto := model.Fatto{
				Anno: anno,
			}
			for _, col := range cmd.Columns {
				switch {
				case col.Name == "CID":
					fatto.CID = row[col.Offset]
				case col.Name == "Codice Individuale":
					fatto.CodiceIndividuale = row[col.Offset]
				case col.Name == "Nome":
					fatto.Nome = row[col.Offset]
				case col.Name == "Cognome":
					fatto.Cognome = row[col.Offset]
				case col.Name == "CognomeNome":
					cognome, nome, ok := strings.Cut(row[col.Offset], " ")
					if ok {
						fatto.Cognome = cognome
						fatto.Nome = nome
					}
				case col.Name == "NomeCognome":
					nome, cognome, ok := strings.Cut(row[col.Offset], " ")
					if ok {
						fatto.Cognome = cognome
						fatto.Nome = nome
					}
				case col.Name == "Dipartimento":
					fatto.Dipartimento = row[col.Offset]
				case col.Name == "Servizio":
					fatto.Servizio = row[col.Offset]
				case col.Name == "Divisione":
					fatto.Divisione = row[col.Offset]
				case col.Name == "Settore":
					fatto.Settore = row[col.Offset]
				case col.Name == "Segmento":
					var s model.Segmento
					err := s.UnmarshalText([]byte(row[col.Offset]))
					if err != nil {
						slog.Error("failed to parse segment", "row", i, "field", col.Name, "segment", row[col.Offset], "error", err)
						continue outer
					}
					fatto.Segmento = s
				case col.Name == "Decorrenza Segmento":
					t, err := safeParseDate(row[col.Offset])
					if err != nil {
						slog.Error("failed to parse date", "row", i, "field", col.Name, "date", row[col.Offset], "error", err)
						return err
					}
					fatto.DecorrenzaSegmento = t
				case col.Name == "Livello":

					match := livello.FindStringSubmatch(row[col.Offset])
					if len(match) > 1 {
						lv, err := strconv.Atoi(match[1])
						if err != nil {
							slog.Error("failed to parse integer", "row", i, "field", col.Name, "integer", row[col.Offset], "error", err)
							return err
						}
						fatto.Livello = lv
					}
				case col.Name == "Decorrenza Livello":
					t, err := safeParseDate(row[col.Offset])
					if err != nil {
						slog.Error("failed to parse date", "row", i, "field", col.Name, "date", row[col.Offset], "error", err)
						return err
					}
					fatto.DecorrenzaLivello = t
				}
			}

			match := false
			if len(cmd.Filters) == 0 {
				match = true
			} else {
			filters:
				for _, filter := range cmd.Filters {
					switch filter.Field {
					case "CID":
						if fatto.CID == filter.Value {
							match = true
							break filters
						}
					case "Codice Individuale":
						if fatto.CodiceIndividuale == filter.Value {
							match = true
							break filters
						}
					case "Nome":
						if fatto.Nome == filter.Value {
							match = true
							break filters
						}
					case "Cognome":
						if fatto.Cognome == filter.Value {
							match = true
							break filters
						}
					case "Dipartimento":
						if fatto.Dipartimento == filter.Value {
							match = true
							break filters
						}
					case "Servizio":
						if fatto.Servizio == filter.Value {
							match = true
							break filters
						}
					case "Divisione":
						if fatto.Divisione == filter.Value {
							match = true
							break filters
						}
					case "Categoria":
						if fatto.Settore == filter.Value {
							match = true
							break filters
						}
					case "Segmento":
						if fatto.Segmento.String() == filter.Value {
							match = true
							break filters
						}
					case "Decorrenza Segmento":
						if fatto.DecorrenzaLivello.Format("02/01/2006") == filter.Value {
							match = true
							break filters
						}
					case "Livello":
						value, err := strconv.Atoi(filter.Value)
						if err != nil {
							slog.Error("failed to parse integer", "integer", filter.Value, "error", err)
							return err
						}
						if fatto.Livello == value {
							match = true
							break filters
						}
					case "Decorrenza Livello":
						if fatto.DecorrenzaLivello.Format("02/01/2006") == filter.Value {
							match = true
							break filters
						}
					}
				}
			}

			e, err := encoder.New(cmd.Format, encoder.WithIndentation(), encoder.WithDataMapper(model.MapFattoToSlice))
			if err != nil {
				slog.Error("failed to create encoder", "format", cmd.Format, "error", err)
				return err
			}
			defer e.Close()

			if match {
				if err := e.Encode(os.Stdout, fatto); err != nil {
					slog.Error("failed to encode fatto", "error", err)
					return err
				}
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
			}

			if _, err := db.NamedExec(`INSERT INTO fatti (anno, cid, codice_individuale, nome, cognome, dipartimento, servizio, divisione, settore, segmento, decorrenza_segmento, livello, decorrenza_livello) VALUES (:anno, :cid, :codice_individuale, :nome, :cognome, :dipartimento, :servizio, :divisione, :settore, :segmento, :decorrenza_segmento, :livello, :decorrenza_livello)`, &fatto); err != nil {
				slog.Error("failed to insert fatto", "fatto", fatto, "error", err)
				return err
				// } else {
				// 	fmt.Printf(".")
			}
		}
	}
	fmt.Println("")

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
