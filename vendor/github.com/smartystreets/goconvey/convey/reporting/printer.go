package reporting

import (
	"fmt"
	"io"
	"strings"
)

type Printer struct {
	out    io.Writer
	prefix string
}

func (self *Printer) Println(message string, values ...any) {
	formatted := self.format(message, values...) + newline
	self.out.Write([]byte(formatted))
}

func (self *Printer) Print(message string, values ...any) {
	formatted := self.format(message, values...)
	self.out.Write([]byte(formatted))
}

func (self *Printer) Insert(text string) {
	self.out.Write([]byte(text))
}

func (self *Printer) format(message string, values ...any) string {
	var formatted string
	if len(values) == 0 {
		formatted = self.prefix + message
	} else {
		formatted = self.prefix + fmt_Sprintf(message, values...)
	}
	indented := strings.Replace(formatted, newline, newline+self.prefix, -1)
	return strings.TrimRight(indented, space)
}

// Extracting fmt.Sprintf to a separate variable circumvents go vet, which, as of go 1.10 is run with go test.
var fmt_Sprintf = fmt.Sprintf

func (self *Printer) Indent() {
	self.prefix += pad
}

func (self *Printer) Dedent() {
	if len(self.prefix) >= padLength {
		self.prefix = self.prefix[:len(self.prefix)-padLength]
	}
}

func NewPrinter(out io.Writer) *Printer {
	self := new(Printer)
	self.out = out
	return self
}

const space = " "
const pad = space + space
const padLength = len(pad)
