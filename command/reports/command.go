package reports

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/excel/collections"
	"github.com/dihedron/excel/logic"
	"github.com/dihedron/excel/model"
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

	var fatti []model.Fatto
	err = db.Select(&fatti, `SELECT * from fatti where dipartimento = 'DIPARTIMENTO INFORMATICA' order by cid asc, ANNO asc`)
	if err != nil {
		slog.Error("failed to query database", "error", err)
		return err
	}

	// {
	// 	PersonaleSDDC, PersonaleRSA, PersonaleDipIT, _ := logic.PersonalePerAnno(fatti)
	// 	EleggibiliConsigliereSDDC, EleggibiliConsigliereRSA, EleggibiliConsigliereDipIT, _ := logic.EleggibiliAConsiglierePerAnno(fatti)
	// 	EleggibiliDirettoreSDDC, EleggibiliDirettoreRSA, EleggibiliDirettoreDipIT, _ := logic.EleggibiliADirettorePerAnno(fatti)
	// 	anni := collections.Merge(collections.MapKeys(PersonaleSDDC), collections.MapKeys(PersonaleRSA), collections.MapKeys(PersonaleDipIT))

	// 	fmt.Printf("A+--------+----------+--------+------------+-------------+-----------+-------------+-----+\n")
	// 	fmt.Printf("A|  Anno  | Strutt.  |  Tot.  |  El. cons. | Prom. Cons. |  El. Dir. |  Prom. Dir. |  %%  |\n")
	// 	fmt.Printf("A+--------+----------+--------+------------+-------------+-----------+-------------+-----+\n")
	// 	for _, anno := range anni {
	// 		fmt.Printf("B|  %4d  |  %-6s  | %6d | %10d | %11d | %9d |\n",
	// 			anno,
	// 			"SDDC",
	// 			len(PersonaleSDDC[anno]),
	// 			len(EleggibiliConsigliereSDDC[anno]),
	// 			-1,
	// 			len(EleggibiliDirettoreSDDC[anno]),
	// 		)
	// 		fmt.Printf("C+--------+----------+--------+------------+-------------+-----------+-----------+\n")
	// 		fmt.Printf("C|  %4d  |  %-6s  | %6d | %10d | %11d | %9d |\n",
	// 			anno,
	// 			"RSA",
	// 			len(PersonaleRSA[anno]),
	// 			len(EleggibiliConsigliereRSA[anno]),
	// 			-1,
	// 			len(EleggibiliDirettoreRSA[anno]),
	// 		)
	// 		fmt.Printf("D+--------+----------+--------+------------+-------------+-----------+-----------+\n")
	// 		fmt.Printf("D|  %4d  |  %-6s  | %6d | %10d | %11d | %9d |\n",
	// 			anno,
	// 			"Dip.",
	// 			len(PersonaleDipIT[anno]),
	// 			len(EleggibiliConsigliereDipIT[anno]),
	// 			-1,
	// 			len(EleggibiliDirettoreDipIT[anno]),
	// 		)
	// 		fmt.Printf("E+--------+----------+--------+------------+-------------+-----------+-----------+\n")
	// 	}
	// }

	{
		EspertiSDDC, EspertiRSA, EspertiDipIT, _ := logic.EspertiPerAnno(fatti)
		EleggibiliConsigliereSDDC, EleggibiliConsigliereRSA, EleggibiliConsigliereDipIT, _ := logic.EleggibiliAConsiglierePerAnno(fatti)
		PromossiConsigliereSDDC, PromossiConsigliereRSA, PromossiConsigliereDipIT, _ := logic.PromossiAConsiglierePerAnno(fatti)
		//EleggibiliDirettoreSDDC, EleggibiliDirettoreRSA, EleggibiliDirettoreDipIT, _ := logic.EleggibiliADirettorePerAnno(fatti)
		anni := collections.Merge(collections.MapKeys(EspertiSDDC), collections.MapKeys(EspertiRSA), collections.MapKeys(EspertiDipIT))

		fmt.Printf("A+--------+---------+-------------+--------+------------+----------+---------+\n")
		fmt.Printf("A|  Anno  |  Unità  | Avanzam. a  | Totale | Eleggibili | Promossi |    %%    |\n")
		fmt.Printf("A+--------+---------+-------------+--------+------------+----------+---------+\n")
		for _, anno := range anni {
			fmt.Printf("B|  %4d  |  %-5s  | Consigliere | %6d | %10d | %8d |  %3.1f%% |\n",
				anno,
				"SDDC",
				len(EspertiSDDC[anno])+len(PromossiConsigliereSDDC),
				len(EleggibiliConsigliereSDDC[anno])+len(PromossiConsigliereSDDC),
				len(PromossiConsigliereSDDC[anno]),
				safePercent(len(PromossiConsigliereSDDC), len(EleggibiliConsigliereSDDC[anno])+len(PromossiConsigliereSDDC)),
			)
			fmt.Printf("C+--------+---------+-------------+--------+------------+----------+--------+\n")
			fmt.Printf("C|  %4d  |  %-5s  | Consigliere | %6d | %10d | %8d |  %3.1f%% |\n",
				anno,
				"RSA",
				len(EspertiRSA[anno])+len(PromossiConsigliereRSA),
				len(EleggibiliConsigliereRSA[anno])+len(PromossiConsigliereRSA),
				len(PromossiConsigliereRSA[anno]),
				safePercent(len(PromossiConsigliereRSA), len(EleggibiliConsigliereRSA[anno])+len(PromossiConsigliereRSA)),
			)
			fmt.Printf("D+--------+---------+-------------+--------+------------+----------+--------+\n")
			fmt.Printf("D|  %4d  |  %-5s  | Consigliere | %6d | %10d | %8d |  %3.1f%% |\n",
				anno,
				"Dpt.",
				len(EspertiDipIT[anno])+len(PromossiConsigliereDipIT),
				len(EleggibiliConsigliereDipIT[anno])+len(PromossiConsigliereDipIT),
				len(PromossiConsigliereDipIT[anno]),
				safePercent(len(PromossiConsigliereDipIT), len(EleggibiliConsigliereDipIT[anno])+len(PromossiConsigliereDipIT)),
			)
			fmt.Printf("E+--------+---------+-------------+--------+------------+----------+--------+\n")
		}
	}

	// {

	// 	fmt.Println("Personale dell'Area Manageriale:")
	// 	for anno, personale := range collections.SortedMap(sddc) {
	// 		fmt.Printf("SDDC[%d]: %d\n", anno, len(personale))
	// 	}
	// 	for anno, personale := range collections.SortedMap(rsa) {
	// 		fmt.Printf("RSA[%d]: %d\n", anno, len(personale))
	// 	}
	// 	for anno, personale := range collections.SortedMap(dipartimento) {
	// 		fmt.Printf("DIP[%d]: %d\n", anno, len(personale))
	// 	}
	// }

	/*
		err = logic.Avanzamenti(db)
		if err != nil {
			slog.Error("failed to generate reports", "error", err)
			return err
		}
	*/

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

func safePercent(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) * 100 / float64(b)
}
