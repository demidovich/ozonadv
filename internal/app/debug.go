package app

import (
	"io"

	"github.com/fatih/color"
)

type Debug struct {
	out io.Writer
}

func newDebug(out io.Writer) Debug {
	return Debug{out: out}
}

func (d Debug) Println(m ...any) {
	_, _ = color.RGB(70, 70, 70).Fprintln(d.out, m...)
}

func (d Debug) Printf(format string, m ...any) {
	_, _ = color.RGB(70, 70, 70).Fprintf(d.out, format, m...)
}
