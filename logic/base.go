package logic

import (
	"fmt"

	"github.com/dihedron/excel/model"
)

func PersonalePerAnno(fatti []model.Fatto) (map[int][]model.Fatto, map[int][]model.Fatto, map[int][]model.Fatto, error) {

	sddc := map[int][]model.Fatto{}
	rsa := map[int][]model.Fatto{}
	dipartimento := map[int][]model.Fatto{}

	for _, fatto := range fatti {

		// popola la cache dei CID
		cidToName.Put(fatto.CID, fmt.Sprintf("%s %s", fatto.Cognome, fatto.Nome))

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
	return sddc, rsa, dipartimento, nil
}

func ConsiglieriPerAnno(fatti []model.Fatto) (map[int][]model.Fatto, map[int][]model.Fatto, map[int][]model.Fatto, error) {

	sddc := map[int][]model.Fatto{}
	rsa := map[int][]model.Fatto{}
	dipartimento := map[int][]model.Fatto{}

	for _, fatto := range fatti {

		// popola la cache dei CID
		cidToName.Put(fatto.CID, fmt.Sprintf("%s %s", fatto.Cognome, fatto.Nome))

		if fatto.Segmento != model.Consigliere {
			continue
		}

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
	return sddc, rsa, dipartimento, nil
}

func EspertiPerAnno(fatti []model.Fatto) (map[int][]model.Fatto, map[int][]model.Fatto, map[int][]model.Fatto, error) {

	sddc := map[int][]model.Fatto{}
	rsa := map[int][]model.Fatto{}
	dipartimento := map[int][]model.Fatto{}

	for _, fatto := range fatti {

		// popola la cache dei CID
		cidToName.Put(fatto.CID, fmt.Sprintf("%s %s", fatto.Cognome, fatto.Nome))

		if fatto.Segmento != model.Esperto {
			continue
		}

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
	return sddc, rsa, dipartimento, nil
}
