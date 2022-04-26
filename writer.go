package kail

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

var prefixColors = []*color.Color{
	color.New(color.FgRed, color.Bold),
	color.New(color.FgGreen, color.Bold),
	color.New(color.FgYellow, color.Bold),
	color.New(color.FgBlue, color.Bold),
	color.New(color.FgMagenta, color.Bold),
	color.New(color.FgCyan, color.Bold),
}

var formatter = &colorjson.Formatter{
	KeyColor:        color.New(color.FgWhite, color.Italic),
	StringColor:     color.New(color.FgHiWhite),
	BoolColor:       color.New(color.FgHiGreen),
	NumberColor:     color.New(color.FgHiCyan),
	NullColor:       color.New(color.FgHiMagenta),
	StringMaxLength: 0,
	DisabledColor:   false,
	Indent:          0,
	RawStrings:      false,
}

var prefixMap = make(map[string]int)

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
		val = (len(prefixMap) + 1) % len(prefixColors)
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
