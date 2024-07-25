package printers

import (
	"jolang2/src/jo"
)

type Printer interface {
	PrintUnit(unit *jo.Unit) string
	Filename(unit *jo.Unit) string
}
