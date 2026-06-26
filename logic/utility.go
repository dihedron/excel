package logic

import (
	"strings"

	"github.com/dihedron/excel/model"
)

func isInSDDC(fatto model.Fatto) bool {
	return strings.Contains(fatto.Settore, "SOFTWARE DEFINED")
}

func isInRSA(fatto model.Fatto) bool {
	return strings.Contains(fatto.Divisione, "RETI, SICUREZZA")
}
