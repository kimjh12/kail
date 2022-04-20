package kail

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var prefixColors = []*color.Color{
	color.New(color.FgRed, color.Bold),
	color.New(color.FgGreen, color.Bold),
	color.New(color.FgYellow, color.Bold),
	color.New(color.FgBlue, color.Bold),
	color.New(color.FgMagenta, color.Bold),
	color.New(color.FgCyan, color.Bold),
	color.New(color.FgWhite, color.Bold),
}

var prefixMap = make(map[string]int)

type Writer interface {
	Print(event Event) error
	Fprint(w io.Writer, event Event) error
}

func NewWriter(out io.Writer) Writer {
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

	if _, err := out.Write(log); err != nil {
		return err
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
