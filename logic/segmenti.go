package logic

import (
	"github.com/dihedron/excel/model"
)

func EleggibiliAConsiglierePerAnno(fatti []model.Fatto) (map[int][]model.Fatto, map[int][]model.Fatto, map[int][]model.Fatto, error) {

	sddc := map[int][]model.Fatto{}
	rsa := map[int][]model.Fatto{}
	dipartimento := map[int][]model.Fatto{}

	for _, fatto := range fatti {

		// popola la cache dei CID
		cidToName.Put(fatto.CID, fatto.Nominativo)

		if fatto.Segmento == model.Esperto && fatto.Livello >= 3 {
			if isInSDDC(fatto) {
				if _, ok := sddc[fatto.Anno]; !ok {
					sddc[fatto.Anno] = []model.Fatto{fatto}
				} else {
					sddc[fatto.Anno] = append(sddc[fatto.Anno], fatto)
				}
			} else if isInRSA(fatto) {
				if _, ok := rsa[fatto.Anno]; !ok {
					rsa[fatto.Anno] = []model.Fatto{fatto}
				} else {
					rsa[fatto.Anno] = append(rsa[fatto.Anno], fatto)
				}
			} else {
				if _, ok := dipartimento[fatto.Anno]; !ok {
					dipartimento[fatto.Anno] = []model.Fatto{fatto}
				} else {
					dipartimento[fatto.Anno] = append(dipartimento[fatto.Anno], fatto)
				}
			}
		}
	}
	return sddc, rsa, dipartimento, nil
}

func EleggibiliADirettorePerAnno(fatti []model.Fatto) (map[int][]model.Fatto, map[int][]model.Fatto, map[int][]model.Fatto, error) {

	sddc := map[int][]model.Fatto{}
	rsa := map[int][]model.Fatto{}
	dipartimento := map[int][]model.Fatto{}

	for _, fatto := range fatti {

		// popola la cache dei CID
		cidToName.Put(fatto.CID, fatto.Nominativo)

		if fatto.Segmento == model.Consigliere && fatto.Livello >= 7 {
			if isInSDDC(fatto) {
				if _, ok := sddc[fatto.Anno]; !ok {
					sddc[fatto.Anno] = []model.Fatto{fatto}
				} else {
					sddc[fatto.Anno] = append(sddc[fatto.Anno], fatto)
				}
			} else if isInRSA(fatto) {
				if _, ok := rsa[fatto.Anno]; !ok {
					rsa[fatto.Anno] = []model.Fatto{fatto}
				} else {
					rsa[fatto.Anno] = append(rsa[fatto.Anno], fatto)
				}
			} else {
				if _, ok := dipartimento[fatto.Anno]; !ok {
					dipartimento[fatto.Anno] = []model.Fatto{fatto}
				} else {
					dipartimento[fatto.Anno] = append(dipartimento[fatto.Anno], fatto)
				}
			}
		}
	}
	return sddc, rsa, dipartimento, nil
}

func PromossiAConsiglierePerAnno(fatti []model.Fatto) (map[int][]PassaggioDiSegmento, map[int][]PassaggioDiSegmento, map[int][]PassaggioDiSegmento, error) {

	sddc := map[int][]PassaggioDiSegmento{}
	rsa := map[int][]PassaggioDiSegmento{}
	dipartimento := map[int][]PassaggioDiSegmento{}

	cid := ""
	segmento := model.Segmento(3)
	livello := 0
	for _, fatto := range fatti {

		// popola la cache dei CID
		cidToName.Put(fatto.CID, fatto.Nominativo)

		if cid == "" {
			cid = fatto.CID
		}

		// seleziona i passaggi di segmento: il segmento diminuisce rispetto all'anno precedente
		// e anche il livello diminuisce (per il reinquadramento nel nuovo segmento)
		if segmento == 2 && fatto.Segmento == 1 && livello > fatto.Livello {

			//fmt.Printf("%s %s passa da %v a %v", fatto.Cognome, fatto.Nome, segmento, fatto.Segmento)

			passaggio := PassaggioDiSegmento{
				CID:        cid,
				Anno:       fatto.Anno,
				Nominativo: fatto.Nominativo,
				Precedente: segmento,
				Attuale:    fatto.Segmento,
				Decorrenza: fatto.DecorrenzaSegmento,
			}

			if isInSDDC(fatto) {
				if _, ok := sddc[fatto.Anno]; !ok {
					sddc[fatto.Anno] = []PassaggioDiSegmento{passaggio}
				} else {
					sddc[fatto.Anno] = append(sddc[fatto.Anno], passaggio)
				}
			} else if isInRSA(fatto) {
				if _, ok := rsa[fatto.Anno]; !ok {
					rsa[fatto.Anno] = []PassaggioDiSegmento{passaggio}
				} else {
					rsa[fatto.Anno] = append(rsa[fatto.Anno], passaggio)
				}
			} else {
				if _, ok := dipartimento[fatto.Anno]; !ok {
					dipartimento[fatto.Anno] = []PassaggioDiSegmento{passaggio}
				} else {
					dipartimento[fatto.Anno] = append(dipartimento[fatto.Anno], passaggio)
				}
			}
		}
		segmento = fatto.Segmento
		livello = fatto.Livello
	}
	return sddc, rsa, dipartimento, nil
}

/*
		if fatto.Segmento == model.Esperto && fatto.Livello >= 3 {
			if isInSDDC(fatto) {
				if _, ok := sddc[fatto.Anno]; !ok {
					sddc[fatto.Anno] = []model.Fatto{fatto}
				} else {
					sddc[fatto.Anno] = append(sddc[fatto.Anno], fatto)
				}
			} else if isInRSA(fatto) {
				if _, ok := rsa[fatto.Anno]; !ok {
					rsa[fatto.Anno] = []model.Fatto{fatto}
				} else {
					rsa[fatto.Anno] = append(rsa[fatto.Anno], fatto)
				}
			} else {
				if _, ok := dipartimento[fatto.Anno]; !ok {
					dipartimento[fatto.Anno] = []model.Fatto{fatto}
				} else {
					dipartimento[fatto.Anno] = append(dipartimento[fatto.Anno], fatto)
				}
			}
		}
	}
	return sddc, rsa, dipartimento, nil
}
*/
