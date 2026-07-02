package encoder

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

type Encoder interface {
	Encode(w io.Writer, data any) error
	Close() error
}

type Option func(e Encoder) error

func New(format string, opts ...Option) (Encoder, error) {
	var e Encoder
	switch format {
	case "json":
		e = &JSONEncoder{}
	case "yaml":
		e = &YAMLEncoder{}
	case "text":
		e = &TextEncoder{}
	case "csv":
		e = &CSVEncoder{}
	case "none":
		e = &NullEncoder{}
	default:
		return nil, errors.New("unsupported format")
	}
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

type JSONEncoder struct {
	indent bool
}

func WithIndentation() Option {
	return func(e Encoder) error {
		if je, ok := e.(*JSONEncoder); ok {
			je.indent = true
		}
		return nil
	}
}

func (e *JSONEncoder) Encode(w io.Writer, data any) error {
	var (
		result []byte
		err    error
	)
	if e.indent {
		result, err = json.MarshalIndent(data, "", "  ")
	} else {
		result, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", string(result))
	return nil
}

func (e *JSONEncoder) Close() error {
	return nil
}

type YAMLEncoder struct{}

func (e *YAMLEncoder) Encode(w io.Writer, data any) error {
	result, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "---\n%s", string(result))
	return nil
}

func (e *YAMLEncoder) Close() error {
	return nil
}

type TextEncoder struct{}

func (e *TextEncoder) Encode(w io.Writer, data any) error {
	fmt.Fprintf(w, "%+v\n", data)
	return nil
}

func (e *TextEncoder) Close() error {
	return nil
}

type CSVEncoder struct {
	writer *csv.Writer
	mapper func(data any) ([]string, error)
}

func WithDataMapper(mapper func(data any) ([]string, error)) Option {
	return func(e Encoder) error {
		if ce, ok := e.(*CSVEncoder); ok {
			ce.mapper = mapper
		}
		return nil
	}
}

func (e *CSVEncoder) Encode(w io.Writer, data any) error {
	if e.writer == nil {
		e.writer = csv.NewWriter(w)
	}

	var (
		record []string
		err    error
	)
	if e.mapper != nil {
		record, err = e.mapper(data)
		if err != nil {
			return err
		}
	} else {
		return errors.New("no mapper function provided")
	}
	e.writer.Write(record)
	return nil
}

func (e *CSVEncoder) Close() error {
	e.writer.Flush()
	return nil
}

type NullEncoder struct{}

func (e *NullEncoder) Encode(w io.Writer, data any) error {
	return nil
}

func (e *NullEncoder) Close() error {
	return nil
}
