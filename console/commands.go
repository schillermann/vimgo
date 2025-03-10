package console

import (
	"fmt"
	"io"
	"os"
)

const (
	escape = "\033"
	csi    = escape + "[" // control sequence introducer
	osc    = escape + "]" // operating system command
)

type Commands struct {
	writer io.ReadWriter
}

func NewCommands() *Commands {
	return &Commands{
		writer: os.Stdout,
	}
}

func (self *Commands) CursorGet() (rowNumber int, columnNumber int) {
	fmt.Fprint(self.writer, csi+"6n")
	fmt.Fscanf(self.writer, csi+"%d;%dR", &rowNumber, &columnNumber)
	return rowNumber, columnNumber
}

func (self *Commands) ColorInverse() {
	fmt.Fprint(self.writer, csi+"7m")
}

func (self *Commands) CursorHide() {
	fmt.Fprint(self.writer, csi+"?25l")
}

func (self *Commands) CursorJumpDown(jump int) {
	fmt.Fprintf(self.writer, csi+"%dB", jump)
}

func (self *Commands) CursorJumpLeft(jump int) {
	fmt.Fprintf(self.writer, csi+"%dD", jump)
}

func (self *Commands) CursorJumpRight(jump int) {
	fmt.Fprintf(self.writer, csi+"%dC", jump)
}

func (self *Commands) CursorJumpUp(jump int) {
	fmt.Fprintf(self.writer, csi+"%dA", jump)
}

func (self *Commands) CursorPositionRestore() {
	fmt.Fprint(self.writer, escape+"8")
}

func (self *Commands) CursorPositionSave() {
	fmt.Fprint(self.writer, escape+"7")
}

func (self *Commands) CursorSet(row int, column int) {
	fmt.Fprintf(self.writer, csi+"%d;%dH", row, column)
}

func (self *Commands) CursorSetTopLeft() {
	fmt.Fprint(self.writer, csi+"H")
}

func (self *Commands) CursorShow() {
	fmt.Fprint(self.writer, csi+"?25h")
}

func (self *Commands) CursorStyleLine() {
	fmt.Fprint(self.writer, csi+"6 q")
}

func (self *Commands) CursorStyleReset() {
	fmt.Fprint(self.writer, csi+"0 q")
}

func (self *Commands) Reset() {
	fmt.Fprint(self.writer, escape+"c")
}

func (self *Commands) ResetFormatting() {
	fmt.Fprint(self.writer, csi+"0m")
}

func (self *Commands) RunePrint(row int, column int, char rune) {
	fmt.Fprintf(self.writer, csi+"%d;%dH%c", row, column, char)
}

func (self *Commands) TitleSet(title string) {
	fmt.Fprintf(self.writer, osc+"0;%s\007", title)
}
