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
	self.TitleSet(self.file.Name())
	return nil
}

func (self *Editor) FileSave() error {
	err := self.file.Save()
	self.StatuslineRender()
	return err
}

func (self *Editor) CursorCheckColumnEnd() {
	columnEnd := len(self.file.Rows()[self.cursorRow-1]) + 1
	if self.cursorColumn > columnEnd {
		self.cursorColumn = columnEnd
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
	if self.cursorColumn > len(self.file.Rows()[self.cursorRow-1]) {
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

func (self *Editor) RowRender(row int) error {
	_, consoleWindowColumns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	self.consoleCommands.CursorHide()
	for columnIndex, char := range self.file.Rows()[row-1] {
		if columnIndex >= consoleWindowColumns {
			self.consoleCommands.CursorShow()
			return nil
		}
		self.consoleCommands.RunePrint(row, columnIndex+1, char)
	}
	self.RowRenderSpaces(row, len(self.file.Rows()[row-1])+1, consoleWindowColumns)
	self.consoleCommands.CursorShow()
	return nil
}

func (self *Editor) RowRenderSpaces(row int, columnStart int, columnEnd int) {
	self.consoleCommands.RunePrint(row, columnStart, '↲')
	for column := columnStart + 1; column <= columnEnd; column++ {
		self.consoleCommands.RunePrint(row, column, ' ')
	}
}

func (self *Editor) RuneDelete() error {
	self.file.RuneDelete(self.cursorRow, self.cursorColumn)
	self.cursorColumn--
	if err := self.RowRender(self.cursorRow); err != nil {
		return err
	}
	if err := self.StatuslineRender(); err != nil {
		return err
	}
	return nil
}

func (self *Editor) RuneInsert(char rune) error {
	self.file.RuneInsert(self.cursorRow, self.cursorColumn, char)
	self.cursorColumn++
	if err := self.RowRender(self.cursorRow); err != nil {
		return err
	}
	if err := self.StatuslineRender(); err != nil {
		return err
	}
	return nil
}

func (self *Editor) ScreenRender() error {
	self.consoleCommands.ScreenClear()

	if err := self.StatuslineRender(); err != nil {
		return err
	}

	consoleWindowRows, _, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	for row := 1; row < consoleWindowRows; row++ {
		if row > len(self.file.Rows()) {
			self.consoleCommands.RunePrint(row, 1, '~')
			continue
		}

		self.RowRender(row)
	}
	self.consoleCommands.CursorMoveTopLeft()

	return nil
}

func (self *Editor) Statusline(middleColumns int, columnsCut int) string {
	fileModified := ' '
	if self.file.Modified() {
		fileModified = '*'
	}
	filename := self.file.Name()
	if columnsCut > 0 {
		filename = ".." + filename[columnsCut+2:]
	}

	return fmt.Sprintf(
		" %[1]s: %[2]s%[3]c%[4]sRow %[5]d/%[6]d, Col %[7]d ",
		strings.ToUpper(self.mode),
		filename,
		fileModified,
		strings.Repeat(" ", middleColumns),
		self.cursorRow,
		self.file.NumberOfRows(),
		self.cursorColumn,
	)
}

func (self *Editor) StatuslineRender() error {
	consoleWindowRows, consoleWindowColumns, err := consoleWindow.Size()
	if err != nil {
		return err
	}

	self.consoleCommands.ColorInverse()

	statuslineMiddleColumns := 1
	statusline := self.Statusline(statuslineMiddleColumns, 0)
	statuslineColumns := len(statusline)
	if statuslineColumns < consoleWindowColumns {
		statusline = self.Statusline(consoleWindowColumns+statuslineMiddleColumns-statuslineColumns, 0)
	} else if statuslineColumns > consoleWindowColumns {
		statusline = self.Statusline(statuslineMiddleColumns, statuslineColumns-consoleWindowColumns)
	}

	self.consoleCommands.CursorHide()
	for index, char := range statusline {
		self.consoleCommands.RunePrint(consoleWindowRows, index+1, char)
	}
	self.consoleCommands.CursorShow()
	self.consoleCommands.Reset()
	self.consoleCommands.CursorMoveTo(self.cursorRow, self.cursorColumn)

	return nil
}

func (self *Editor) TitleSet(title string) {
	self.consoleCommands.TitleSet("VimGo - " + title)
}
