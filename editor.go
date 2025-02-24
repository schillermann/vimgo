package main

import (
	"fmt"
	"strings"

	"github.com/schillermann/vimgo/console"
)

type Editor struct {
	cursorRow       int
	cursorColumn    int
	consoleCommands *console.Commands
	file            *File
	mode            string
}

func NewEditor(file *File) *Editor {
	return &Editor{
		cursorRow:       1,
		cursorColumn:    1,
		consoleCommands: console.NewCommands(),
		file:            file,
		mode:            "view",
	}
}

func (self *Editor) FileLoad() error {
	if err := self.file.Load(); err != nil {
		return err
	}
	self.consoleCommands.CursorMoveTopLeft()
	self.TitleSet(self.file.Filename())
	return nil
}

func (self *Editor) FileSave() error {
	err := self.file.Save()
	self.StatuslineRender()
	return err
}

func (self *Editor) CursorCheckColumnEnd() {
	columnLength := len(self.file.Rows()[self.cursorRow-1])
	if self.cursorColumn > columnLength {
		self.cursorColumn = columnLength
		self.consoleCommands.CursorMoveTo(self.cursorRow, self.cursorColumn)
	}
}

func (self *Editor) CursorMoveLeft(jump int) {
	if self.cursorColumn < 2 {
		return
	}
	self.consoleCommands.CursorMoveLeft(1)
	self.cursorColumn--
	self.StatuslineRender()
}

func (self *Editor) CursorMoveDown(jump int) {
	if self.cursorRow >= len(self.file.Rows()) {
		return
	}
	self.consoleCommands.CursorMoveDown(1)
	self.cursorRow++
	self.CursorCheckColumnEnd()
	self.StatuslineRender()
}

func (self *Editor) CursorMoveUp(jump int) {
	if self.cursorRow < 2 {
		return
	}
	self.consoleCommands.CursorMoveUp(1)
	self.cursorRow--
	self.CursorCheckColumnEnd()
	self.StatuslineRender()
}

func (self *Editor) CursorMoveRight(jump int) {
	if self.cursorColumn >= len(self.file.Rows()[self.cursorRow-1]) {
		return
	}
	self.consoleCommands.CursorMoveRight(1)
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

func (self *Editor) ModeToEdit() {
	self.mode = "edit"
	self.StatuslineRender()
}

func (self *Editor) ModeToView() {
	self.mode = "view"
	self.StatuslineRender()
}

func (self *Editor) RowRender(rowNumber int) {
	for columnIndex, char := range self.file.Rows()[rowNumber-1] {
		self.consoleCommands.RunePrint(rowNumber, columnIndex+1, char)
	}
}

func (self *Editor) RuneInsert(char rune) {
	self.file.Insert(self.cursorRow, self.cursorColumn, char)
	self.RowRender(self.cursorRow)
	self.cursorColumn++
	self.consoleCommands.CursorMoveTo(self.cursorRow, self.cursorColumn)
	self.StatuslineRender()
}

func (self *Editor) ScreenRender() error {
	self.consoleCommands.ScreenClear()

	if err := self.StatuslineRender(); err != nil {
		return err
	}

	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	for rowIndex := 0; rowIndex < rows-1; rowIndex++ {
		if rowIndex >= len(self.file.Rows()) {
			self.consoleCommands.RunePrint(rowIndex+1, 1, '~')
			continue
		}

		for columnIndex, char := range self.file.Rows()[rowIndex] {
			if columnIndex >= columns {
				break
			}
			self.consoleCommands.RunePrint(rowIndex+1, columnIndex+1, char)
		}
	}
	self.consoleCommands.CursorMoveTopLeft()

	return nil
}

func (self *Editor) StatuslineRender() error {
	rows, columns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	self.consoleCommands.ColorInverse()

	fileModified := ' '
	if self.file.Modified() {
		fileModified = '*'
	}
	leftStatusline := fmt.Sprintf(" %s: %c%s", strings.ToUpper(self.mode), fileModified, self.file.Filename())
	rightStatusline := fmt.Sprintf("Row %d/%d, Col %d ", self.cursorRow, self.file.NumberOfRows(), self.cursorColumn)
	middleStatuslineLength := columns - len(leftStatusline) - len(rightStatusline)
	statusline := leftStatusline + strings.Repeat(" ", middleStatuslineLength) + rightStatusline

	self.consoleCommands.CursorHide()
	for index, char := range statusline {
		self.consoleCommands.RunePrint(rows, index+1, char)
	}
	self.consoleCommands.CursorShow()
	self.consoleCommands.Reset()
	self.consoleCommands.CursorMoveTo(self.cursorRow, self.cursorColumn)

	return nil
}

func (self *Editor) TitleSet(title string) {
	self.consoleCommands.TitleSet("VimGo - " + title)
}
