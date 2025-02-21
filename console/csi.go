package console

import (
	"fmt"
	"io"
)

const (
	escape                  = "\033"
	csi                     = escape + "[" // control sequence introducer
	AnsiCustorMovePrintRune = csi + "%d;%dH%c"
)

type Csi struct {
	Writer io.Writer
}

func (self *Csi) ClearScreen() {
	self.MoveCursorLeftCorner()
	fmt.Fprint(self.Writer, csi+"2J")
}

func (self *Csi) MoveCursorLeftCorner() {
	fmt.Fprint(self.Writer, csi+"H")
}

func (self *Csi) MoveCursorDown(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dB", jump)
}

func (self *Csi) MoveCursorLeft(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dD", jump)
}

func (self *Csi) MoveCursorRight(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dC", jump)
}

func (self *Csi) MoveCursorUp(jump int) {
	fmt.Fprintf(self.Writer, csi+"%dA", jump)
}

func (self *Csi) PrintRune(row int, column int, char rune) {
	fmt.Fprintf(self.Writer, csi+"%d;%dH%c", row, column, char)
}
