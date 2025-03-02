package main

import (
	"github.com/schillermann/vimgo/console"
)

type Editor struct {
	file               *File
	fileCursor         *FileCursor
	editorMode         *EditorMode
	statusline         *Statusline
	consoleCommands    *console.Commands
	screenCursorRow    int
	screenCursorColumn int
	columnsWidth       int
	rowsHeight         int
}

func NewEditor(
	file *File,
	fileCursor *FileCursor,
	editorMode *EditorMode,
	statusline *Statusline,
	consoleCommands *console.Commands,
	rowsHeight int,
	columnsWidth int,
) *Editor {
	return &Editor{
		file:               file,
		fileCursor:         fileCursor,
		editorMode:         editorMode,
		statusline:         statusline,
		consoleCommands:    consoleCommands,
		screenCursorRow:    1,
		screenCursorColumn: 1,
		columnsWidth:       columnsWidth,
		rowsHeight:         rowsHeight,
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
	self.statusline.Render()
	return err
}

func (self *Editor) CursorMoveLeft(jump int) {
	self.fileCursor.MoveLeft(jump)
	self.statusline.Render()
}

func (self *Editor) CursorMoveDown(jump int) {
	self.fileCursor.MoveDown(jump)
	self.statusline.Render()
}

func (self *Editor) CursorMoveUp(jump int) {
	self.fileCursor.MoveUp(jump)
	self.statusline.Render()
}

func (self *Editor) CursorMoveRight(jump int) {
	self.fileCursor.MoveRight(jump)
	self.statusline.Render()
}

func (self *Editor) ModeToEdit() {
	self.editorMode.ToEdit()
	self.statusline.Render()
}

func (self *Editor) ModeToView() {
	self.editorMode.ToView()
	self.statusline.Render()
}

func (self *Editor) RowRender(row int) {
	self.consoleCommands.CursorHide()
	for columnIndex, char := range self.file.Rows()[row-1] {
		if columnIndex >= self.columnsWidth {
			self.consoleCommands.CursorShow()
		}
		self.consoleCommands.RunePrint(row, columnIndex+1, char)
	}
	columnEnd := len(self.file.Rows()[row-1])
	self.consoleCommands.RunePrint(row, columnEnd+1, '↲')
	self.RowRenderSpaces(row, columnEnd+2, self.columnsWidth)
	self.consoleCommands.CursorShow()
}

func (self *Editor) RowRenderSpaces(row int, columnStart int, columnEnd int) {
	for column := columnStart; column <= columnEnd; column++ {
		self.consoleCommands.RunePrint(row, column, ' ')
	}
}

func (self *Editor) RuneDelete() {
	self.file.RuneDelete(self.fileCursor.GetRow(), self.fileCursor.GetColumn())
	self.fileCursor.MoveLeft(1)
	self.RowRender(self.fileCursor.GetRow())
	self.statusline.Render()
}

func (self *Editor) RuneInsert(char rune) {
	self.file.RuneInsert(self.fileCursor.GetRow(), self.fileCursor.GetColumn(), char)
	self.fileCursor.MoveRight(1)
	self.RowRender(self.fileCursor.GetRow())
	self.statusline.Render()
}

func (self *Editor) ScreenRender() {
	self.statusline.Render()

	for row := 1; row <= self.rowsHeight; row++ {
		if row > len(self.file.Rows()) {
			self.consoleCommands.RunePrint(row, 1, '~')
			self.RowRenderSpaces(row, 2, self.columnsWidth)
			continue
		}

		self.RowRender(row)
	}
	self.consoleCommands.CursorMoveTopLeft()
}

func (self *Editor) SetColumnsWidth(columns int) {
	if self.columnsWidth == columns {
		return
	}
	self.columnsWidth = columns
	self.ScreenRender()
}

func (self *Editor) SetRowsHeight(rows int) {
	if self.rowsHeight == rows {
		return
	}
	self.rowsHeight = rows
	self.ScreenRender()
}

func (self *Editor) TitleSet(title string) {
	self.consoleCommands.TitleSet("VimGo - " + title)
}
