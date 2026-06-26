package load

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/dihedron/excel/model"
)

type Converter func(string) (any, error)

func ToInt() Converter {
	return func(s string) (any, error) {
		return strconv.Atoi(s)
	}
}

func ToTime(format string) Converter {
	return func(s string) (any, error) {
		return time.Parse(format, s)
	}
}

func ToSegmento() Converter {
	return func(value string) (any, error) {

		var s model.Segmento
		err := s.UnmarshalText([]byte(value))
		if err != nil {
			slog.Error("failed to parse segment", "value", value, "error", err)
			return nil, err
		}
		return s, nil
	}
}

// func convertInt(s string) (any, error) {
// 	return strconv.Atoi(s)
// }

// func convertString(s string) (any, error) {
// 	return s, nil
// }
