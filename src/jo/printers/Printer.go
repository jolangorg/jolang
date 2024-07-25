package printers

import (
	"github.com/jolangorg/jolang/src/jo"
)

type Printer interface {
	PrintUnit(unit *jo.Unit) string
	Filename(unit *jo.Unit) string
}
