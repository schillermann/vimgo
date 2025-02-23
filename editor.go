package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/schillermann/vimgo/console"
)

type Editor struct {
	cursorRow    int
	cursorColumn int
	consoleCsi   *console.Csi
	file         *File
	mode         string
}

func NewEditor(file *File) *Editor {
	return &Editor{
		cursorRow:    1,
		cursorColumn: 1,
		consoleCsi:   console.NewCsi(),
		file:         file,
		mode:         "view",
	}
}

func (self *Editor) FileLoad() error {
	if err := self.file.Load(); err != nil {
		return err
	}
	self.consoleCsi.CursorMoveTopLeft()
	return nil
}

func (self *Editor) CursorCheckColumnEnd() {
	columnLength := len(self.file.Rows()[self.cursorRow-1])
	if self.cursorColumn > columnLength {
		self.cursorColumn = columnLength
		self.consoleCsi.CursorMoveTo(self.cursorRow, self.cursorColumn)
	}
}

func (self *Editor) CursorMoveLeft(jump int) {
	if self.cursorColumn < 2 {
		return
	}
	self.consoleCsi.CursorMoveLeft(1)
	self.cursorColumn--
	self.StatuslineRender()
}

func (self *Editor) CursorMoveDown(jump int) {
	if self.cursorRow >= len(self.file.Rows()) {
		return
	}
	self.consoleCsi.CursorMoveDown(1)
	self.cursorRow++
	self.CursorCheckColumnEnd()
	self.StatuslineRender()
}

func (self *Editor) CursorMoveUp(jump int) {
	if self.cursorRow < 2 {
		return
	}
	self.consoleCsi.CursorMoveUp(1)
	self.cursorRow--
	self.CursorCheckColumnEnd()
	self.StatuslineRender()
}

func (self *Editor) CursorMoveRight(jump int) {
	if self.cursorColumn >= len(self.file.Rows()[self.cursorRow-1]) {
		return
	}
	self.consoleCsi.CursorMoveRight(1)
	self.cursorColumn++
	self.StatuslineRender()
}

func (self *Editor) IsModeEdit() bool {
	if self.mode == "edit" {
		return true
	}
	return false
}
func (self *Editor) IsModeView() bool {
	if self.mode == "view" {
		return true
	}
	return false
}

func (self *Editor) RowRender(rowNumber int) {
	for columnIndex, char := range self.file.Rows()[rowNumber-1] {
		self.consoleCsi.RunePrint(rowNumber, columnIndex+1, char)
	}
}

func (self *Editor) RuneInsert(char rune) {
	self.file.Rows()[self.cursorRow-1] = slices.Insert(self.file.Rows()[self.cursorRow-1], self.cursorColumn-1, char)
	self.RowRender(self.cursorRow)
	self.cursorColumn++
	self.consoleCsi.CursorMoveTo(self.cursorRow, self.cursorColumn)
	self.StatuslineRender()
}

func (self *Editor) ScreenRender() error {
	self.consoleCsi.ScreenClear()

	if err := self.StatuslineRender(); err != nil {
		return err
	}

	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	for rowIndex := 0; rowIndex < rows-1; rowIndex++ {
		if rowIndex >= len(self.file.Rows()) {
			self.consoleCsi.RunePrint(rowIndex+1, 1, '~')
			continue
		}

		for columnIndex, char := range self.file.Rows()[rowIndex] {
			if columnIndex >= columns {
				break
			}
			self.consoleCsi.RunePrint(rowIndex+1, columnIndex+1, char)
		}
	}
	self.consoleCsi.CursorMoveTopLeft()

	return nil
}

func (self *Editor) StatuslineRender() error {
	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	self.consoleCsi.ColorInverse()

	fileModified := ' '
	leftStatusline := fmt.Sprintf(" %s: %c%s", strings.ToUpper(self.mode), fileModified, self.file.Filename())
	rightStatusline := fmt.Sprintf("Row %d/%d, Col %d ", self.cursorRow, self.file.NumberOfRows(), self.cursorColumn)
	middleStatuslineLength := columns - len(leftStatusline) - len(rightStatusline)
	statusline := leftStatusline + strings.Repeat(" ", middleStatuslineLength) + rightStatusline

	self.consoleCsi.CursorHide()
	for index, char := range statusline {
		self.consoleCsi.RunePrint(rows, index+1, char)
	}
	self.consoleCsi.CursorShow()
	self.consoleCsi.Reset()
	self.consoleCsi.CursorMoveTo(self.cursorRow, self.cursorColumn)

	return nil
}

func (self *Editor) ModeToEdit() {
	self.mode = "edit"
	self.StatuslineRender()
}

func (self *Editor) ModeToView() {
	self.mode = "view"
	self.StatuslineRender()
}
