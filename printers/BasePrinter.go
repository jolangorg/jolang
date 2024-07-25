package printers

import "fmt"

type BasePrinter struct {
	Buffer            string
	Indent            int
	ShouldPrintIndent bool
}

func NewBasePrinter() *BasePrinter {
	return &BasePrinter{
		Buffer: "",
		Indent: 0,
	}
}

func (p *BasePrinter) Write(b []byte) (n int, err error) {
	p.Buffer += string(b)
	return len(b), nil
}

func (p *BasePrinter) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(p, format, a...)
}

func (p *BasePrinter) Print(strs ...string) {
	if p.ShouldPrintIndent {
		p.PrintIndent()
		p.ShouldPrintIndent = false
	}
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

func (u *BasePrinter) PrintIndent() {
	for i := 0; i < u.Indent; i++ {
		u.Buffer += "\t"
	}
}
