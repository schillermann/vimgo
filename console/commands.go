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
	writer io.Writer
}

func NewCommands() *Commands {
	return &Commands{
		writer: os.Stdout,
	}
}

func (self *Commands) ColorInverse() {
	fmt.Fprint(self.writer, csi+"7m")
}

func (self *Commands) CursorHide() {
	fmt.Fprint(self.writer, csi+"?25l")
}

func (self *Commands) CursorMoveDown(jump int) {
	fmt.Fprintf(self.writer, csi+"%dB", jump)
}

func (self *Commands) CursorMoveLeft(jump int) {
	fmt.Fprintf(self.writer, csi+"%dD", jump)
}

func (self *Commands) CursorMoveRight(jump int) {
	fmt.Fprintf(self.writer, csi+"%dC", jump)
}

func (self *Commands) CursorMoveTo(row int, column int) {
	fmt.Fprintf(self.writer, csi+"%d;%dH", row, column)
}

func (self *Commands) CursorMoveTopLeft() {
	fmt.Fprint(self.writer, csi+"H")
}

func (self *Commands) CursorMoveUp(jump int) {
	fmt.Fprintf(self.writer, csi+"%dA", jump)
}

func (self *Commands) CursorShow() {
	fmt.Fprint(self.writer, csi+"?25h")
}

func (self *Commands) Reset() {
	fmt.Fprint(self.writer, csi+"0m")
}

func (self *Commands) RunePrint(row int, column int, char rune) {
	fmt.Fprintf(self.writer, csi+"%d;%dH%c", row, column, char)
}

func (self *Commands) ScreenClear() {
	self.CursorMoveTopLeft()
	fmt.Fprint(self.writer, csi+"0J")
}

func (self *Commands) TitleSet(title string) {
	fmt.Fprintf(self.writer, osc+"0;%s\007", title)
}
