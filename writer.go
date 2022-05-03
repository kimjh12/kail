package kail

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

var prefixColors = []*color.Color{
	color.New(color.FgHiRed),
	color.New(color.FgHiGreen),
	color.New(color.FgHiYellow),
	color.New(color.FgHiBlue),
	color.New(color.FgHiMagenta),
	color.New(color.FgHiCyan),
}

var formatter = &colorjson.Formatter{
	KeyColor:        color.New(color.FgWhite, color.Faint),
	StringColor:     color.New(color.FgHiWhite),
	BoolColor:       color.New(color.FgHiGreen),
	NumberColor:     color.New(color.FgHiCyan),
	NullColor:       color.New(color.FgHiMagenta),
	StringMaxLength: 0,
	DisabledColor:   false,
	Indent:          0,
	RawStrings:      false,
}

const prefixMapMaxSize = 50

var prefixMap = make(map[string]int, prefixMapMaxSize)

type Writer interface {
	Print(event Event) error
	Fprint(w io.Writer, event Event) error
}

func NewWriter(out io.Writer, indent int) Writer {
	formatter.Indent = indent
	return &writer{out}
}

type writer struct {
	out io.Writer
}

func (w *writer) Print(ev Event) error {
	return w.Fprint(w.out, ev)
}

func (w *writer) Fprint(out io.Writer, ev Event) error {
	prefix := w.prefix(ev)

	var prefixColor *color.Color
	if val, ok := prefixMap[prefix]; ok {
		prefixColor = prefixColors[val]
	} else {
		if len(prefixMap) >= prefixMapMaxSize {
			prefixMap = make(map[string]int, prefixMapMaxSize)
		}
		val = len(prefixMap) % len(prefixColors)
		prefixMap[prefix] = val
		prefixColor = prefixColors[val]
	}

	if _, err := prefixColor.Fprint(out, prefix); err != nil {
		return err
	}
	if _, err := prefixColor.Fprint(out, ": "); err != nil {
		return err
	}

	log := ev.Log()

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(log), &obj); err != nil {
		if _, err := out.Write(log); err != nil {
			return err
		}
	} else {
		s, _ := formatter.Marshal(obj)
		if _, err := out.Write(s); err != nil {
			return err
		}
	}

	if sz := len(log); sz == 0 || log[sz-1] != byte('\n') {
		if _, err := out.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) prefix(ev Event) string {
	return fmt.Sprintf("%v/%v[%v]",
		ev.Source().Namespace(),
		ev.Source().Name(),
		ev.Source().Container())
}
