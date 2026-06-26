package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Fatto struct {
	ID                 int64     `db:"id" json:"-" yaml:"-"`
	Anno               int       `db:"anno" json:"anno,omitempty" yaml:"anno,omitempty"`
	CID                string    `db:"cid" json:"cid,omitempty" yaml:"cid,omitempty"`
	CodiceIndividuale  string    `db:"codice_individuale" json:"codice_individuale,omitempty" yaml:"codice_individuale,omitempty"`
	Nominativo         string    `db:"nominativo" json:"nominativo,omitempty" yaml:"nominativo,omitempty"`
	Dipartimento       string    `db:"dipartimento" json:"dipartimento,omitempty" yaml:"dipartimento,omitempty"`
	Servizio           string    `db:"servizio" json:"servizio,omitempty" yaml:"servizio,omitempty"`
	Divisione          string    `db:"divisione" json:"divisione,omitempty" yaml:"divisione,omitempty"`
	Settore            string    `db:"settore" json:"settore,omitempty" yaml:"settore,omitempty"`
	Segmento           Segmento  `db:"segmento" json:"segmento" yaml:"segmento"`
	DecorrenzaSegmento time.Time `db:"decorrenza_segmento" json:"decorrenza_segmento,omitempty" yaml:"decorrenza_segmento,omitempty"`
	Livello            int       `db:"livello" json:"livello" yaml:"livello"`
	DecorrenzaLivello  time.Time `db:"decorrenza_livello" json:"decorrenza_livello,omitempty" yaml:"decorrenza_livello,omitempty"`
}

const (
	SchemaFatti = `
CREATE TABLE IF NOT EXISTS fatti (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	anno INTEGER,
	cid TEXT,
	codice_individuale TEXT,
	nominativo TEXT,
	dipartimento TEXT,
	servizio TEXT,
	divisione TEXT,
	settore TEXT,
	segmento TEXT,
	decorrenza_segmento DATE,
	livello INTEGER,
	decorrenza_livello DATE
);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_fatti_cid_anno ON fatti(cid, anno);
`
)

func MapFattoToSlice(data any) ([]string, error) {
	f, ok := data.(Fatto)
	if !ok {
		return nil, errors.New("failed to cast to Fatto")
	}
	return []string{fmt.Sprintf("%d", f.Anno), f.CID, f.CodiceIndividuale, f.Nominativo, f.Dipartimento, f.Servizio, f.Divisione, f.Settore, f.Segmento.String(), f.DecorrenzaSegmento.Format("02/01/2006"), strconv.Itoa(f.Livello), f.DecorrenzaLivello.Format("02/01/2006")}, nil
}
