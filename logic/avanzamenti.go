package logic

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/excel/model"
	"github.com/jmoiron/sqlx"
)

type Avanzamento struct {
	ID           string
	Anno         int
	Cognome      string
	Nome         string
	FromSegmento int
	ToSegmento   int
	Decorrenza   time.Time
	NelSettore   bool
}

func Avanzamenti(db *sqlx.DB) error {

	// Connect using sqlx (wraps standard database/sql)
	db, err := sqlx.Connect("sqlite3", "excel.db")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	defer db.Close()

	var cids []string
	err = db.Select(&cids, `SELECT distinct(CID) from fatti where DIPARTIMENTO = 'DIPARTIMENTO INFORMATICA' and Anno = 2026	`)
	if err != nil {
		slog.Error("failed to query database", "error", err)
		return err
	}

	avanzamentiBySegmento := map[int][]Avanzamento{}

	for _, cid := range cids {

		var facts []model.Fatto
		err = db.Select(&facts, `SELECT * from fatti where CID = ? and dipartimento = 'DIPARTIMENTO INFORMATICA' order by ANNO asc`, cid)
		if err != nil {
			slog.Error("failed to query database", "error", err)
			return err
		}

		segmento := 3
		for _, fact := range facts {
			if int(fact.Segmento) < segmento && segmento < 3 {
				avanzamento := Avanzamento{
					ID:           cid,
					Anno:         fact.Anno,
					Cognome:      fact.Cognome,
					Nome:         fact.Nome,
					FromSegmento: segmento,
					ToSegmento:   int(fact.Segmento),
					Decorrenza:   fact.DecorrenzaSegmento,
				}
				//if fact.CID
				avanzamentiBySegmento[segmento] = append(avanzamentiBySegmento[segmento], avanzamento)
			}
			segmento = int(fact.Segmento)
		}
	}

	for segmento, avanzamenti := range avanzamentiBySegmento {
		fmt.Printf("Segmento: %d\n", segmento)
		fmt.Printf("Numero di avanzamenti: %d\n", len(avanzamenti))
		for _, avanzamento := range avanzamenti {
			fmt.Printf("%v\n", avanzamento)
		}
	}

	return nil
}
