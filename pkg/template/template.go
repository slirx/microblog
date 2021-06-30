package template

import (
	"fmt"
	html "html/template"
	"path"
	"strings"
	text "text/template"

	"github.com/pkg/errors"
)

type GeneratorType int

const (
	TypeHTML GeneratorType = iota + 1
	TypeText
)

type Generator interface {
	Generate(t GeneratorType, fileName string, data interface{}) (string, error)
}

type generator struct {
}

func (g generator) Generate(gType GeneratorType, fileName string, data interface{}) (string, error) {
	var builder strings.Builder

	switch gType {
	case TypeHTML:
		t, err := html.New(path.Base(fileName)).ParseFiles(fileName)
		if err != nil {
			return "", errors.WithStack(err)
		}

		if err = t.Execute(&builder, data); err != nil {
			return "", errors.WithStack(err)
		}

		return builder.String(), nil
	case TypeText:
		t, err := text.New(path.Base(fileName)).ParseFiles(fileName)
		if err != nil {
			return "", errors.WithStack(err)
		}

		if err = t.Execute(&builder, data); err != nil {
			return "", errors.WithStack(err)
		}

		return builder.String(), nil
	}

	return "", errors.WithStack(fmt.Errorf("invalid type: %d", gType))
}

func New() Generator {
	return generator{}
}
