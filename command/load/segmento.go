package load

import (
	"encoding/json"
	"fmt"

	"go.yaml.in/yaml/v3"
)

type Segmento int

const (
	Direttore Segmento = iota
	Consigliere
	Esperto
)

func (s *Segmento) String() string {
	if s == nil {
		return ""
	}
	return []string{"Direttore", "Consigliere", "Esperto"}[int(*s)]
}

func (s Segmento) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Segmento) UnmarshalText(text []byte) error {
	seg := string(text)
	switch seg {
	case "Direttore":
		*s = Direttore
	case "Consigliere":
		*s = Consigliere
	case "Esperto":
		*s = Esperto
	default:
		return fmt.Errorf("invalid segment: %s", seg)
	}
	return nil
}

func (s Segmento) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Segmento) UnmarshalJSON(value []byte) error {
	return s.UnmarshalText(value)
}

func (s Segmento) MarshalYAML() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Segmento) UnmarshalYAML(value *yaml.Node) error {
	return s.UnmarshalText([]byte(value.Value))
}
