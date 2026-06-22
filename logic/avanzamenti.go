package logic

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/dihedron/excel/model"
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

const (
	patternSettoreSDDC  = "SOFTWARE DEFINED"
	patternDivisioneRSA = "RETI, SICUREZZA"
)

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

	personeInDivisioneRSAPerAnno := map[int][]string{
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

	personeNelDipartimentoPerAnno := map[int][]string{
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

	eleggibiliInDivisioneRSAPerAnnoESegmento := map[int]map[model.Segmento][]string{
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

	// accumulatori globali per gli avanzamenti di segmento
	totaleEleggibiliAlSegmentoInSDDC := 0
	totaleEleggibiliAlSegmentoInRSA := 0
	totaleEleggibiliAlSegmentoNelDipartimento := 0
	totalePromossiAlSegmentoInSDDC := 0
	totalePromossiAlSegmentoInRSA := 0
	totalePromossiAlSegmentoNelDipartimento := 0

	// accumulatori globali per gli aumenti al primo anno
	totaleAumentiAlPrimoAnnoInSDDC := 0
	totaleAumentiAlPrimoAnnoInRSA := 0
	totaleAumentiAlPrimoAnnoNelDipartimento := 0

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
			if strings.Contains(fatto.Settore, patternSettoreSDDC) {
				eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento], fatto.CID)
				totaleEleggibiliAlSegmentoInSDDC++
			} else if strings.Contains(fatto.Divisione, patternDivisioneRSA) {
				eleggibiliInDivisioneRSAPerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliInDivisioneRSAPerAnnoESegmento[fatto.Anno][fatto.Segmento], fatto.CID)
				totaleEleggibiliAlSegmentoInRSA++
			} else {
				eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
				totaleEleggibiliAlSegmentoNelDipartimento++
			}

		}

		// seleziona gli eleggibili per gli avanzamenti di segmento a direttore (minimo livello 7)
		if fatto.Segmento == 1 && fatto.Livello >= 7 {
			if strings.Contains(fatto.Settore, patternSettoreSDDC) {
				eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelSettorePerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
				totaleEleggibiliAlSegmentoInSDDC++
			} else if strings.Contains(fatto.Divisione, patternDivisioneRSA) {
				eleggibiliInDivisioneRSAPerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliInDivisioneRSAPerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
				totaleEleggibiliAlSegmentoInRSA++
			} else {
				eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento] = append(eleggibiliNelDipartimentoPerAnnoESegmento[fatto.Anno][fatto.Segmento], cid)
				totaleEleggibiliAlSegmentoNelDipartimento++
			}
		}

		// estrai le persone che in quello specifico anno erano nel settore SDDC o in RSA
		if strings.Contains(fatto.Settore, patternSettoreSDDC) {
			personeNelSettorePerAnno[fatto.Anno] = append(personeNelSettorePerAnno[fatto.Anno], cid)
		} else if strings.Contains(fatto.Divisione, patternDivisioneRSA) {
			personeInDivisioneRSAPerAnno[fatto.Anno] = append(personeInDivisioneRSAPerAnno[fatto.Anno], cid)
		} else {
			personeNelDipartimentoPerAnno[fatto.Anno] = append(personeNelDipartimentoPerAnno[fatto.Anno], cid)
		}

		// seleziona i passaggi di segmento: il segmento diminuisce rispetto all'anno precedente
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

			// conta i promossi al segmento
			if strings.Contains(fatto.Settore, patternSettoreSDDC) {
				totalePromossiAlSegmentoInSDDC++
			} else if strings.Contains(fatto.Divisione, patternDivisioneRSA) {
				totalePromossiAlSegmentoInRSA++
			} else {
				totalePromossiAlSegmentoNelDipartimento++
			}

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

			// conta gli aumenti al primo anno
			if aumento.AlPrimoAnno {
				if strings.Contains(fatto.Settore, patternSettoreSDDC) {
					totaleAumentiAlPrimoAnnoInSDDC++
				} else if strings.Contains(fatto.Divisione, patternDivisioneRSA) {
					totaleAumentiAlPrimoAnnoInRSA++
				} else {
					totaleAumentiAlPrimoAnnoNelDipartimento++
				}
			}

			// aggiorna il livello e l'anno con quello di questo avanzamento
			livello = fatto.Livello
			anno = fatto.Anno
		}
	}

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
			fmt.Printf("   - Persone nella Divisione RSA: %d\n", len(personeInDivisioneRSAPerAnno[anno]))

			// calcola le persone eleggibili per la promozione in SDDC
			eleggibiliInSDDC := []string{}
			for _, cid := range eleggibiliNelSettorePerAnnoESegmento[anno][model.Segmento(segmento+1)] {
				eleggibiliInSDDC = append(eleggibiliInSDDC, cidToName[cid])
			}
			totaleEleggibiliAlSegmentoInSDDC += len(eleggibiliInSDDC)

			eleggibiliInRSA := []string{}
			for _, cid := range eleggibiliInDivisioneRSAPerAnnoESegmento[anno][model.Segmento(segmento+1)] {
				eleggibiliInRSA = append(eleggibiliInRSA, cidToName[cid])
			}
			totaleEleggibiliAlSegmentoInRSA += len(eleggibiliInRSA)

			eleggibiliNelDipartimento := []string{}
			for _, cid := range eleggibiliNelDipartimentoPerAnnoESegmento[anno][model.Segmento(segmento+1)] {
				eleggibiliNelDipartimento = append(eleggibiliNelDipartimento, cidToName[cid])
			}
			totaleEleggibiliAlSegmentoNelDipartimento += len(eleggibiliNelDipartimento)

			// nel caso ci sia stato un passggio di segmento, il promosso non è più tra gli eleggibili
			// quindi dobbiamo riaggiungerlo (in testa!) se lo vediamo ormai promosso
			for _, passaggio := range passaggi {
				if isInSet(passaggio.CID, personeNelSettorePerAnno[anno]) {
					eleggibiliInSDDC = append([]string{cidToName[passaggio.CID]}, eleggibiliInSDDC...)
				}
				if isInSet(passaggio.CID, personeInDivisioneRSAPerAnno[anno]) {
					eleggibiliInRSA = append([]string{cidToName[passaggio.CID]}, eleggibiliInRSA...)
				}
			}
			fmt.Printf("   - Eleggibili nel Settore SDDC: %d (%v)\n", len(eleggibiliInSDDC), strings.Join(eleggibiliInSDDC, ", "))
			fmt.Printf("   - Eleggibili nella Divisione RSA: %d (%v)\n", len(eleggibiliInRSA), strings.Join(eleggibiliInRSA, ", "))
			fmt.Printf("   - Eleggibili nel Dipartimento: %d\n", len(eleggibiliNelDipartimentoPerAnnoESegmento[anno][model.Segmento(segmento+1)]))
			fmt.Printf("   - Numero di avanzamenti: %d\n", len(passaggi))
			passaggiInSDDC := 0
			passaggiInRSA := 0
			passaggiNelDipartimento := 0
			for _, passaggio := range passaggi {
				if isInSet(passaggio.CID, personeNelSettorePerAnno[anno]) {
					fmt.Printf("     - %s %s (%s)\n", passaggio.Cognome, passaggio.Nome, passaggio.CID)
					passaggiInSDDC++
				} else if isInSet(passaggio.CID, personeInDivisioneRSAPerAnno[anno]) {
					fmt.Printf("     - %s %s (%s)\n", passaggio.Cognome, passaggio.Nome, passaggio.CID)
					passaggiInRSA++
				} else {
					fmt.Printf("     - %s %s (%s)\n", passaggio.Cognome, passaggio.Nome, passaggio.CID)
					passaggiNelDipartimento++
				}
			}
			fmt.Printf("   - %% di promossi in SDDC: %.1f%% (%d/%d)\n", safePercent(passaggiInSDDC, len(eleggibiliInSDDC)), passaggiInSDDC, len(eleggibiliInSDDC))
			fmt.Printf("   - %% di promossi in RSA: %.1f%% (%d/%d)\n", safePercent(passaggiInRSA, len(eleggibiliInRSA)), passaggiInRSA, len(eleggibiliInRSA))
			fmt.Printf("   - %% di promossi nel Dipartimento: %.1f%% (%d/%d)\n", safePercent(passaggiNelDipartimento, len(eleggibiliNelDipartimentoPerAnnoESegmento[anno][model.Segmento(segmento+1)])), passaggiNelDipartimento, len(eleggibiliNelDipartimentoPerAnnoESegmento[anno][model.Segmento(segmento+1)]))
		}

		// aumenti di livello
		aumentiPerAnno := aumenti[anno]
		fmt.Printf("----------------------------------------------------\n")
		fmt.Printf(" - Aumenti di livello:\n")
		numeroAumentiAlPrimoAnno := 0
		numeroAumentiAlPrimoAnnoInSDDC := []string{}
		numeroAumentiAlPrimoAnnoInRSA := []string{}
		numeroAumentiNonAlPrimoAnno := 0
		numeroAumentiNonAlPrimoAnnoInSDDC := []string{}
		numeroAumentiNonAlPrimoAnnoInRSA := []string{}
		for _, aumento := range aumentiPerAnno {
			if aumento.AlPrimoAnno {
				numeroAumentiAlPrimoAnno++
				if isInSet(aumento.CID, personeNelSettorePerAnno[anno]) {
					numeroAumentiAlPrimoAnnoInSDDC = append(numeroAumentiAlPrimoAnnoInSDDC, fmt.Sprintf("%s %s", aumento.Cognome, aumento.Nome))
				} else if isInSet(aumento.CID, personeInDivisioneRSAPerAnno[anno]) {
					numeroAumentiAlPrimoAnnoInRSA = append(numeroAumentiAlPrimoAnnoInRSA, fmt.Sprintf("%s %s", aumento.Cognome, aumento.Nome))
				}
			} else {
				numeroAumentiNonAlPrimoAnno++
				if isInSet(aumento.CID, personeNelSettorePerAnno[anno]) {
					numeroAumentiNonAlPrimoAnnoInSDDC = append(numeroAumentiNonAlPrimoAnnoInSDDC, fmt.Sprintf("%s %s", aumento.Cognome, aumento.Nome))
				} else if isInSet(aumento.CID, personeInDivisioneRSAPerAnno[anno]) {
					numeroAumentiNonAlPrimoAnnoInRSA = append(numeroAumentiNonAlPrimoAnnoInRSA, fmt.Sprintf("%s %s", aumento.Cognome, aumento.Nome))
				}
			}
		}
		fmt.Printf("   - Numero di aumenti al primo anno: %d\n", numeroAumentiAlPrimoAnno)
		fmt.Printf("       di cui %d su %d nel settore SDDC (%v)\n", len(numeroAumentiAlPrimoAnnoInSDDC), len(personeNelSettorePerAnno[anno]), strings.Join(numeroAumentiAlPrimoAnnoInSDDC, ", "))
		fmt.Printf("       di cui %d su %d nella divisione RSA (%v)\n", len(numeroAumentiAlPrimoAnnoInRSA), len(personeInDivisioneRSAPerAnno[anno]), strings.Join(numeroAumentiAlPrimoAnnoInRSA, ", "))
		fmt.Printf("   - Numero di aumenti non al primo anno: %d\n", numeroAumentiNonAlPrimoAnno)
		fmt.Printf("       di cui %d su %d nel settore SDDC (%v)\n", len(numeroAumentiNonAlPrimoAnnoInSDDC), len(personeNelSettorePerAnno[anno]), strings.Join(numeroAumentiNonAlPrimoAnnoInSDDC, ", "))
		fmt.Printf("       di cui %d su %d nella divisione RSA (%v)\n", len(numeroAumentiNonAlPrimoAnnoInRSA), len(personeInDivisioneRSAPerAnno[anno]), strings.Join(numeroAumentiNonAlPrimoAnnoInRSA, ", "))
	}
	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Percentuale avanzamento di segmento nel settore SDDC: %.1f%% (%d/%d)\n", safePercent(totalePromossiAlSegmentoInSDDC, totaleEleggibiliAlSegmentoInSDDC), totalePromossiAlSegmentoInSDDC, totaleEleggibiliAlSegmentoInSDDC)
	fmt.Printf("Percentuale avanzamento di segmento nella divisione RSA: %.1f%% (%d/%d)\n", safePercent(totalePromossiAlSegmentoInRSA, totaleEleggibiliAlSegmentoInRSA), totalePromossiAlSegmentoInRSA, totaleEleggibiliAlSegmentoInRSA)
	fmt.Printf("Percentuale avanzamento di segmento nel dipartimento: %.1f%% (%d/%d)\n", safePercent(totalePromossiAlSegmentoNelDipartimento, totaleEleggibiliAlSegmentoNelDipartimento), totalePromossiAlSegmentoNelDipartimento, totaleEleggibiliAlSegmentoNelDipartimento)
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

func safePercent(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) * 100 / float64(b)
}
