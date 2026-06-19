package load

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
	"go.yaml.in/yaml/v3"
)

type Load struct {
	File    string    `short:"f" long:"file" description:"The Excel file to load." required:"true"`
	Sheets  []Sheet   `short:"s" long:"sheet" description:"The sheet to load (format: <sheet>:<label>)." required:"true"`
	Format  string    `short:"t" long:"format" description:"The format of the output." optional:"true" default:"text" choice:"text" choice:"json" choice:"yaml" choice:"csv"`
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

type Employee struct {
	Year               string    `db:"year" json:"year,omitempty" yaml:"year,omitempty"`
	CID                string    `db:"cid" json:"cid,omitempty" yaml:"cid,omitempty"`
	CodiceIndividuale  string    `db:"codice_individuale" json:"codice_individuale,omitempty" yaml:"codice_individuale,omitempty"`
	Nome               string    `db:"nome" json:"nome,omitempty" yaml:"nome,omitempty"`
	Cognome            string    `db:"cognome" json:"cognome,omitempty" yaml:"cognome,omitempty"`
	Dipartimento       string    `db:"dipartimento" json:"dipartimento,omitempty" yaml:"dipartimento,omitempty"`
	Servizio           string    `db:"servizio" json:"servizio,omitempty" yaml:"servizio,omitempty"`
	Divisione          string    `db:"divisione" json:"divisione,omitempty" yaml:"divisione,omitempty"`
	Settore            string    `db:"categoria" json:"categoria,omitempty" yaml:"categoria,omitempty"`
	Segmento           Segmento  `db:"segmento" json:"segmento" yaml:"segmento"`
	DecorrenzaSegmento time.Time `db:"decorrenza_segmento" json:"decorrenza_segmento,omitempty" yaml:"decorrenza_segmento,omitempty"`
	Livello            int       `db:"livello" json:"livello" yaml:"livello"`
	DecorrenzaLivello  time.Time `db:"decorrenza_livello" json:"decorrenza_livello,omitempty" yaml:"decorrenza_livello,omitempty"`
}

const schema = `
CREATE TABLE IF NOT EXISTS employee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT
);`

func (cmd *Load) Execute(args []string) error {

	slog.Debug("opening CVS file", "file", cmd.File, "sheets", cmd.Sheets, "mappings", cmd.Columns, "filters", cmd.Filters)

	// 1. Open the Excel file
	f, err := excelize.OpenFile(cmd.File)
	if err != nil {
		slog.Error("failed to open file", "path", cmd.File, "error", err)
		return err
	}
	defer f.Close()

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
			employee := Employee{
				Year: sheet.Label,
			}
			for _, col := range cmd.Columns {
				switch {
				case col.Name == "CID":
					employee.CID = row[col.Offset]
				case col.Name == "Codice Individuale":
					employee.CodiceIndividuale = row[col.Offset]
				case col.Name == "Nome":
					employee.Nome = row[col.Offset]
				case col.Name == "Cognome":
					employee.Cognome = row[col.Offset]
				case col.Name == "CognomeNome":
					slog.Error("parsing CognomeNome", "row", i, "field", col.Name, "value", row[col.Offset])
					cognome, nome, ok := strings.Cut(row[col.Offset], " ")
					if ok {
						employee.Cognome = cognome
						employee.Nome = nome
					}
				case col.Name == "NomeCognome":
					nome, cognome, ok := strings.Cut(row[col.Offset], " ")
					if ok {
						employee.Cognome = cognome
						employee.Nome = nome
					}
				case col.Name == "Dipartimento":
					employee.Dipartimento = row[col.Offset]
				case col.Name == "Servizio":
					employee.Servizio = row[col.Offset]
				case col.Name == "Divisione":
					employee.Divisione = row[col.Offset]
				case col.Name == "Categoria":
					employee.Settore = row[col.Offset]
				case col.Name == "Segmento":
					var s Segmento
					err := s.UnmarshalText([]byte(row[col.Offset]))
					if err != nil {
						slog.Error("failed to parse segment", "row", i, "field", col.Name, "segment", row[col.Offset], "error", err)
						continue outer
					}
					employee.Segmento = s
				case col.Name == "Decorrenza Segmento":
					t, err := safeParseDate(row[col.Offset])
					if err != nil {
						slog.Error("failed to parse date", "row", i, "field", col.Name, "date", row[col.Offset], "error", err)
						return err
					}
					employee.DecorrenzaSegmento = t
				case col.Name == "Livello":

					match := livello.FindStringSubmatch(row[col.Offset])
					if len(match) > 1 {
						lv, err := strconv.Atoi(match[1])
						if err != nil {
							slog.Error("failed to parse integer", "row", i, "field", col.Name, "integer", row[col.Offset], "error", err)
							return err
						}
						employee.Livello = lv
					}
				case col.Name == "Decorrenza Livello":
					t, err := safeParseDate(row[col.Offset])
					if err != nil {
						slog.Error("failed to parse date", "row", i, "field", col.Name, "date", row[col.Offset], "error", err)
						return err
					}
					employee.DecorrenzaLivello = t
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
						if employee.CID == filter.Value {
							match = true
							break filters
						}
					case "Codice Individuale":
						if employee.CodiceIndividuale == filter.Value {
							match = true
							break filters
						}
					case "Nome":
						if employee.Nome == filter.Value {
							match = true
							break filters
						}
					case "Cognome":
						if employee.Cognome == filter.Value {
							match = true
							break filters
						}
					case "Dipartimento":
						if employee.Dipartimento == filter.Value {
							match = true
							break filters
						}
					case "Servizio":
						if employee.Servizio == filter.Value {
							match = true
							break filters
						}
					case "Divisione":
						if employee.Divisione == filter.Value {
							match = true
							break filters
						}
					case "Categoria":
						if employee.Settore == filter.Value {
							match = true
							break filters
						}
					case "Segmento":
						if employee.Segmento.String() == filter.Value {
							match = true
							break filters
						}
					case "Decorrenza Segmento":
						if employee.DecorrenzaLivello.Format("02/01/2006") == filter.Value {
							match = true
							break filters
						}
					case "Livello":
						value, err := strconv.Atoi(filter.Value)
						if err != nil {
							slog.Error("failed to parse integer", "integer", filter.Value, "error", err)
							return err
						}
						if employee.Livello == value {
							match = true
							break filters
						}
					case "Decorrenza Livello":
						if employee.DecorrenzaLivello.Format("02/01/2006") == filter.Value {
							match = true
							break filters
						}
					}
				}
			}

			var w *csv.Writer

			if match {
				switch cmd.Format {
				case "text":
					fmt.Printf("%+v\n", employee)
				case "json":
					b, err := json.MarshalIndent(employee, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(b))
				case "yaml":
					b, err := yaml.Marshal(employee)
					if err != nil {
						return err
					}
					fmt.Println(string(b))
				case "csv":
					if w == nil {
						w = csv.NewWriter(os.Stdout)
						defer w.Flush()
					}
					w.Write([]string{employee.Year, employee.CID, employee.CodiceIndividuale, employee.Nome, employee.Cognome, employee.Dipartimento, employee.Servizio, employee.Divisione, employee.Settore, employee.Segmento.String(), employee.DecorrenzaSegmento.Format("02/01/2006"), strconv.Itoa(employee.Livello), employee.DecorrenzaLivello.Format("02/01/2006")})
					w.Flush()
				}
			}

		}
	}

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
