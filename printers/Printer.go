package printers

import "jolang2"

type Printer interface {
	PrintUnit(unit *jolang2.Unit) string
	Filename(unit *jolang2.Unit) string
}
