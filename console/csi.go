package console

import (
	"fmt"
	"io"
	"os"
)

const (
	escape = "\033"
	csi    = escape + "[" // control sequence introducer
)

type Csi struct {
	writer io.Writer
}

func NewCsi() *Csi {
	return &Csi{
		writer: os.Stdout,
	}
}

func (self *Csi) ColorInverse() {
	fmt.Fprint(self.writer, csi+"7m")
}

func (self *Csi) CursorHide() {
	fmt.Fprint(self.writer, csi+"?25l")
}

func (self *Csi) CursorMoveDown(jump int) {
	fmt.Fprintf(self.writer, csi+"%dB", jump)
}

func (self *Csi) CursorMoveLeft(jump int) {
	fmt.Fprintf(self.writer, csi+"%dD", jump)
}

func (self *Csi) CursorMoveRight(jump int) {
	fmt.Fprintf(self.writer, csi+"%dC", jump)
}

func (self *Csi) CursorMoveTo(row int, column int) {
	fmt.Fprintf(self.writer, csi+"%d;%dH", row, column)
}

func (self *Csi) CursorMoveTopLeft() {
	fmt.Fprint(self.writer, csi+"H")
}

func (self *Csi) CursorMoveUp(jump int) {
	fmt.Fprintf(self.writer, csi+"%dA", jump)
}

func (self *Csi) CursorShow() {
	fmt.Fprint(self.writer, csi+"?25h")
}

func (self *Csi) Reset() {
	fmt.Fprint(self.writer, csi+"0m")
}

func (self *Csi) RunePrint(row int, column int, char rune) {
	fmt.Fprintf(self.writer, csi+"%d;%dH%c", row, column, char)
}

func (self *Csi) ScreenClear() {
	self.CursorMoveTopLeft()
	fmt.Fprint(self.writer, csi+"0J")
}
