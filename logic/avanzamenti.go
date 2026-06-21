package logic

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dihedron/excel/model"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
)

type PassaggioDiSegmento struct {
	CID        string
	Anno       int
	Cognome    string
	Nome       string
	DaSegmento int
	ASegmento  int
	Decorrenza time.Time
	NelSettore bool
}

type Aumento struct {
	CID         string
	Anno        int
	Cognome     string
	Nome        string
	Livello     int
	AlPrimoAnno bool
}

func Avanzamenti(db *sqlx.DB) error {

	cidToName := make(map[string]string)

	personeNelSettorePerAnno := map[int][]string{
		2017: []string{},
		2018: []string{},
		2019: []string{},
		2020: []string{},
		2021: []string{},
		2022: []string{},
		2023: []string{},
		2024: []string{},
		2025: []string{},
		2026: []string{},
	}

	eleggibiliNelDipartimentoPerAnnoESegmento := map[int]map[model.Segmento][]string{
		2017: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2018: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2019: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2020: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2021: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2022: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2023: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2024: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2025: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2026: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
	}

	eleggibiliNelSettorePerAnnoESegmento := map[int]map[model.Segmento][]string{
		2017: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2018: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2019: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2020: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2021: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2022: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2023: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2024: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2025: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
		2026: map[model.Segmento][]string{
			model.Direttore:   []string{},
			model.Consigliere: []string{},
		},
	}

	passaggiDiSegmentoPerAnno := map[int]map[int][]PassaggioDiSegmento{
		2017: map[int][]PassaggioDiSegmento{},
		2018: map[int][]PassaggioDiSegmento{},
		2019: map[int][]PassaggioDiSegmento{},
		2020: map[int][]PassaggioDiSegmento{},
		2021: map[int][]PassaggioDiSegmento{},
		2022: map[int][]PassaggioDiSegmento{},
		2023: map[int][]PassaggioDiSegmento{},
		2024: map[int][]PassaggioDiSegmento{},
		2025: map[int][]PassaggioDiSegmento{},
		2026: map[int][]PassaggioDiSegmento{},
	}

	aumenti := map[int][]Aumento{}

	var fatti []model.Fatto
	err := db.Select(&fatti, `SELECT * from fatti where dipartimento = 'DIPARTIMENTO INFORMATICA' order by cid asc, ANNO asc`)
	if err != nil {
		slog.Error("failed to query database", "error", err)
		return err
	}

	cid := ""
	segmento := 3
	livello := 0
	anno := 0
	for _, fatto := range fatti {

		// popola la cache dei CID
		if _, ok := cidToName[fatto.CID]; !ok {
			cidToName[fatto.CID] = fmt.Sprintf("%s %s", fatto.Cognome, fatto.Nome)
		}

		if cid == "" {
			cid = fatto.CID
		} else if cid != fatto.CID {
			cid = fatto.CID
			segmento = 3
			livello = 0
			anno = 0
		}

		// seleziona gli eleggibili per gli avanzamenti di segmento a consigliere (minimo livello 3)
		if fatto.Segmento == 2 && fatto.Livello >= 3 {
			if strings.Contains(fatto.Settore, "SOFTWARE DEFINED") {
				eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento], fatto.CID)
			} else {
				eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
			}
		}

		// seleziona gli eleggibili per gli avanzamenti di segmento a direttore (minimo livello 7)
		if fatto.Segmento == 1 && fatto.Livello >= 7 {
			if strings.Contains(fatto.Settore, "SOFTWARE DEFINED") {
				eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
			} else {
				eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
			}
		}

		// estrai le persone che in quello specifico anno erano nel settore SDDC
		if strings.Contains(fatto.Settore, "SOFTWARE DEFINED") {
			personeNelSettorePerAnno[fatto.Anno] = append(personeNelSettorePerAnno[fatto.Anno], cid)
		}

		// seleziona i passaggi di segmento: il segmento diminuisce rispetto all'anno precedente (quindi è un avanzamento)
		if int(fatto.Segmento) < segmento && segmento < 3 {
			passaggioDiSegmento := PassaggioDiSegmento{
				CID:        cid,
				Anno:       fatto.Anno,
				Cognome:    fatto.Cognome,
				Nome:       fatto.Nome,
				DaSegmento: segmento,
				ASegmento:  int(fatto.Segmento),
				Decorrenza: fatto.DecorrenzaSegmento,
			}
			passaggiDiSegmentoPerAnno[fatto.Anno][int(fatto.Segmento)] = append(passaggiDiSegmentoPerAnno[fatto.Anno][int(fatto.Segmento)], passaggioDiSegmento)
		}
		// aggiorna il segmento con quello di questo avanzamento
		segmento = int(fatto.Segmento)

		// calcola aumenti di livello
		if livello == 0 {
			livello = fatto.Livello
			anno = fatto.Anno
		} else if fatto.Livello > livello {
			aumento := Aumento{
				CID:         cid,
				Anno:        fatto.Anno,
				Cognome:     fatto.Cognome,
				Nome:        fatto.Nome,
				Livello:     fatto.Livello,
				AlPrimoAnno: fatto.Anno == (anno + 1),
			}
			aumenti[fatto.Anno] = append(aumenti[fatto.Anno], aumento)

			// aggiorna il livello e l'annocon quello di questo avanzamento
			livello = fatto.Livello
			anno = fatto.Anno
		}
	}
	red := color.New(color.FgRed).FprintfFunc()

	// ordina gli anni
	anni := make([]int, 0)
	for anno := range passaggiDiSegmentoPerAnno {
		anni = append(anni, anno)
	}
	sort.Ints(anni)

	// per ogni anno...
	for _, anno := range anni {
		fmt.Printf("----------------------------------------------------\n")
		fmt.Printf("Anno: %d\n", anno)
		// passaggi di segmento
		passaggiPerAnno := passaggiDiSegmentoPerAnno[anno]
		// ordina sui segmenti
		segmenti := make([]int, 0)
		for segmento := range passaggiPerAnno {
			segmenti = append(segmenti, segmento)
		}
		sort.Ints(segmenti)
		for _, segmento := range segmenti {
			passaggi := passaggiPerAnno[segmento]
			fmt.Printf("----------------------------------------------------\n")
			fmt.Printf(" - Passaggio di Segmento a: %v\n", (model.Segmento(segmento)).String())
			fmt.Printf("   - Persone nel Settore SDDC: %d\n", len(personeNelSettorePerAnno[anno]))

			persone := []string{}
			for _, cid := range eleggibiliNelSettorePerAnnoESegmento[anno][model.Segmento(segmento+1)] {
				persone = append(persone, cidToName[cid])
			}
			fmt.Printf("   - Eleggibili nel Settore SDDC: %d (%v)\n", len(eleggibiliNelSettorePerAnnoESegmento[anno][model.Segmento(segmento+1)]), strings.Join(persone, ", "))
			fmt.Printf("   - Eleggibili nel Dipartimento: %d\n", len(eleggibiliNelDipartimentoPerAnnoESegmento[anno][model.Segmento(segmento+1)]))
			fmt.Printf("   - Numero di avanzamenti: %d\n", len(passaggi))
			for _, passaggio := range passaggi {
				if isInSet(passaggio.CID, personeNelSettorePerAnno[anno]) {
					red(os.Stdout, "     - %s %s (%s)\n", passaggio.Cognome, passaggio.Nome, passaggio.CID)
				} else {
					fmt.Printf("     - %s %s (%s)\n", passaggio.Cognome, passaggio.Nome, passaggio.CID)
				}
			}
		}

		// aumenti di livello
		aumentiPerAnno := aumenti[anno]
		fmt.Printf("----------------------------------------------------\n")
		fmt.Printf(" - Aumenti di livello:\n")
		numeroAumentiAlPrimoAnno := 0
		numeroAumentiAlPrimoAnnoNonInSddc := 0
		numeroAumentiNonAlPrimoAnno := 0
		numeroAumentiNonAlPrimoAnnoNonInSddc := 0
		for _, aumento := range aumentiPerAnno {
			if aumento.AlPrimoAnno {
				numeroAumentiAlPrimoAnno++
				if !isInSet(aumento.CID, personeNelSettorePerAnno[anno]) {
					numeroAumentiAlPrimoAnnoNonInSddc++
				}
			} else {
				numeroAumentiNonAlPrimoAnno++
				if !isInSet(aumento.CID, personeNelSettorePerAnno[anno]) {
					numeroAumentiNonAlPrimoAnnoNonInSddc++
				}
			}
		}
		fmt.Printf("   - Numero di aumenti al primo anno: %d\n", numeroAumentiAlPrimoAnno)
		red(os.Stdout, "       di cui %d nel settore SDDC\n", numeroAumentiAlPrimoAnno-numeroAumentiAlPrimoAnnoNonInSddc)
		fmt.Printf("   - Numero di aumenti non al primo anno: %d\n", numeroAumentiNonAlPrimoAnno)
		red(os.Stdout, "       di cui %d nel settore SDDC\n", numeroAumentiNonAlPrimoAnno-numeroAumentiNonAlPrimoAnnoNonInSddc)
	}
	fmt.Printf("----------------------------------------------------\n")
	return nil
}

func isInSet(value string, set []string) bool {
	for _, s := range set {
		if value == s {
			return true
		}
	}
	return false
}
