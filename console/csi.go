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

func (self *Csi) PrintRune(row int, column int, char rune) {
	fmt.Fprintf(self.Writer, csi+"%d;%dH%c", row, column, char)
}
