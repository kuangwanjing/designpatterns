package format

import (
	"errors"
	"format/method"
)

type FormatManager struct {
	FormatFactory map[string]func() method.Format
}

var manager FormatManager

func Register(fn string, fh func() method.Format) {
	if manager.FormatFactory == nil {
		manager.FormatFactory = make(map[string]func() method.Format)
	}
	manager.FormatFactory[fn] = fh
}

func Stringify(fn string, data interface{}) (string, error) {
	f, ok := manager.FormatFactory[fn]
	if !ok {
		return "", errors.New("factory not found")
	}
	format := f()
	return format.Stringify(data), nil
}
