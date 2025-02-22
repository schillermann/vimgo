package console

import (
	"fmt"
	"io"
)

const (
	escape = "\033"
	csi    = escape + "[" // control sequence introducer
)

type Csi struct {
	Writer io.Writer
}

func (self *Csi) ColorInverse() {
	fmt.Fprint(self.Writer, csi+"7m")
}

func (self *Csi) CursorHide() {
	fmt.Fprint(self.Writer, csi+"?25l")
}

func (self *Csi) CursorMoveDown(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dB", jump)
}

func (self *Csi) CursorMoveLeft(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dD", jump)
}

func (self *Csi) CursorMoveRight(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dC", jump)
}

func (self *Csi) CursorMoveTo(row int, column int) {
	fmt.Fprintf(self.Writer, csi+"%d;%dH", row, column)
}

func (self *Csi) CursorMoveTopLeft() {
	fmt.Fprint(self.Writer, csi+"H")
}

func (self *Csi) CursorMoveUp(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dA", jump)
}

func (self *Csi) CursorShow() {
	fmt.Fprint(self.Writer, csi+"?25h")
}

func (self *Csi) Reset() {
	fmt.Fprint(self.Writer, csi+"0m")
}

func (self *Csi) RunePrint(row int, column int, char rune) {
	fmt.Fprintf(self.Writer, csi+"%d;%dH%c", row, column, char)
}

func (self *Csi) ScreenClear() {
	self.CursorMoveTopLeft()
	fmt.Fprint(self.Writer, csi+"0J")
}
