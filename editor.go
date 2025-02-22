package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/schillermann/vimgo/console"
)

type Editor struct {
	file         *File
	cursorRow    int
	cursorColumn int
}

var consoleCsi = console.Csi{Writer: os.Stdout}

func NewEditor(file *File) *Editor {
	consoleCsi.CursorMoveTopLeft()

	return &Editor{
		file:         file,
		cursorRow:    1,
		cursorColumn: 1,
	}
}

func (self *Editor) FileLoad() error {
	if err := self.file.Load(); err != nil {
		return err
	}
	return nil
}

func (self *Editor) CursorMoveLeft(jump int) {
	if self.cursorColumn < 2 {
		return
	}
	consoleCsi.CursorMoveLeft(1)
	self.cursorColumn--
	self.StatuslineRender()
}

func (self *Editor) CursorMoveDown(jump int) {
	if self.cursorRow >= len(self.file.Rows()) {
		return
	}
	consoleCsi.CursorMoveDown(1)
	self.cursorRow++
	self.StatuslineRender()
}

func (self *Editor) CursorMoveUp(jump int) {
	if self.cursorRow < 2 {
		return
	}
	consoleCsi.CursorMoveUp(1)
	self.cursorRow--
	self.StatuslineRender()
}

func (self *Editor) CursorMoveRight(jump int) {
	if self.cursorColumn >= len(self.file.Rows()[self.cursorRow-1]) {
		return
	}
	consoleCsi.CursorMoveRight(1)
	self.cursorColumn++
	self.StatuslineRender()
}

func (self *Editor) ScreenRender() error {
	consoleCsi.ScreenClear()

	if err := self.StatuslineRender(); err != nil {
		return err
	}

	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	for rowIndex := 0; rowIndex < rows-1; rowIndex++ {
		if rowIndex >= len(self.file.Rows()) {
			consoleCsi.RunePrint(rowIndex+1, 1, '~')
			continue
		}

		for columnIndex, char := range self.file.Rows()[rowIndex] {
			if columnIndex >= columns {
				break
			}
			consoleCsi.RunePrint(rowIndex+1, columnIndex+1, char)
		}
	}
	consoleCsi.CursorMoveTopLeft()

	return nil
}

func (self *Editor) StatuslineRender() error {
	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	consoleCsi.ColorInverse()
	mode := "VIEW"
	fileModified := ' '
	leftStatusline := fmt.Sprintf(" %s: %c%s", mode, fileModified, self.file.Filename())
	rightStatusline := fmt.Sprintf("Row %d/%d, Col %d ", self.cursorRow, self.file.NumberOfRows(), self.cursorColumn)
	middleStatuslineLength := columns - len(leftStatusline) - len(rightStatusline)
	statusline := leftStatusline + strings.Repeat(" ", middleStatuslineLength) + rightStatusline

	consoleCsi.CursorHide()
	for index, char := range statusline {
		consoleCsi.RunePrint(rows, index+1, char)
	}
	consoleCsi.CursorShow()
	consoleCsi.Reset()
	consoleCsi.CursorMoveTo(self.cursorRow, self.cursorColumn)

	return nil
}
