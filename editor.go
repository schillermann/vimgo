package main

import (
	"os"

	"github.com/schillermann/vimgo/console"
)

type Editor struct {
	file         File
	cursorRow    int
	cursorColumn int
}

var consoleCsi = console.Csi{Writer: os.Stdout}

func NewEditor(file File) *Editor {
	consoleCsi.MoveCursorLeftCorner()

	return &Editor{
		file:         file,
		cursorRow:    1,
		cursorColumn: 1,
	}
}

func (self *Editor) LoadFile() error {
	if err := self.file.Read(); err != nil {
		return err
	}
	return nil
}

func (self *Editor) MoveCursorLeft(jump int) {
	if self.cursorColumn < 2 {
		return
	}
	consoleCsi.MoveCursorLeft(1)
	self.cursorColumn--
}

func (self *Editor) MoveCursorDown(jump int) {
	if self.cursorRow >= len(self.file.Rows()) {
		return
	}
	consoleCsi.MoveCursorDown(1)
	self.cursorRow++
}

func (self *Editor) MoveCursorUp(jump int) {
	if self.cursorRow < 2 {
		return
	}
	consoleCsi.MoveCursorUp(1)
	self.cursorRow--
}

func (self *Editor) MoveCursorRight(jump int) {
	if self.cursorColumn >= len(self.file.Rows()[self.cursorRow -1]) {
		return
	}
	consoleCsi.MoveCursorRight(1)
	self.cursorColumn++
}

func (self *Editor) RenderScreen() error {
	consoleCsi.ClearScreen()

	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		if rowIndex >= len(self.file.Rows()) {
			consoleCsi.PrintRune(rowIndex+1, 1, '~')
			continue
		}

		for columnIndex, char := range self.file.Rows()[rowIndex] {
			if columnIndex >= columns {
				break
			}
			consoleCsi.PrintRune(rowIndex+1, columnIndex+1, char)
		}
	}
	consoleCsi.MoveCursorLeftCorner()

	return nil
}
