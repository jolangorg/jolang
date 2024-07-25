package printers

import (
	"fmt"
	"jolang2/src/jo"
)

type BasePrinter struct {
	Project           *jo.Project
	Buffer            string
	Indent            int
	ShouldPrintIndent bool
}

func NewBasePrinter(project *jo.Project) *BasePrinter {
	return &BasePrinter{
		Project: project,
		Buffer:  "",
		Indent:  0,
	}
}

func (p *BasePrinter) Write(b []byte) (n int, err error) {
	p.PrintIndent()
	p.Buffer += string(b)
	return len(b), nil
}

func (p *BasePrinter) Printf(format string, a ...any) {
	p.PrintIndent()
	_, _ = fmt.Fprintf(p, format, a...)
}

func (p *BasePrinter) Print(strs ...string) {
	p.PrintIndent()
	for i, s := range strs {
		if i == 0 {
			p.Buffer += s
		} else {
			p.Buffer += " " + s
		}
	}
}

func (p *BasePrinter) Println(s ...string) {
	p.Print(s...)
	p.Print("\n")
	p.ShouldPrintIndent = true
}

func (p *BasePrinter) PrintIndent() {
	if !p.ShouldPrintIndent {
		return
	}

	p.ShouldPrintIndent = false
	for i := 0; i < p.Indent; i++ {
		p.Buffer += "\t"
	}
}
