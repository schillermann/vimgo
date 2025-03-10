package main

import (
	"github.com/schillermann/vimgo/console"
)

type Editor struct {
	file            *File
	fileCursor      *FileCursor
	editorMode      *EditorMode
	statusline      *Statusline
	consoleCommands *console.Commands
	offsetColumn    int
	offsetRow       int
	columnCount     int
	rowCount        int
	scrollOff       int
}

func NewEditor(
	file *File,
	fileCursor *FileCursor,
	editorMode *EditorMode,
	statusline *Statusline,
	consoleCommands *console.Commands,
	rowCount int,
	columnCount int,
) *Editor {
	return &Editor{
		file:            file,
		fileCursor:      fileCursor,
		editorMode:      editorMode,
		statusline:      statusline,
		consoleCommands: consoleCommands,
		rowCount:        rowCount,
		columnCount:     columnCount,
		scrollOff:       5,
	}
}

func (self *Editor) ColumnNumberGet() int {
	return self.fileCursor.ColumnNumberGet() - self.offsetColumn
}

func (self *Editor) CursorJumpLeft(jump int) {
	self.fileCursor.JumpLeft(jump)

	currentColumnNumber := self.fileCursor.ColumnNumberGet() - self.offsetColumn
	isOffset := self.columnCount < self.file.ColumnEnd(self.fileCursor.RowNumberGet())
	isCrossScrollOff := currentColumnNumber <= self.scrollOff
	isScroll := self.offsetColumn > 0
	if isOffset &&
		isCrossScrollOff &&
		isScroll {
		self.offsetColumn--
		self.Render()
	}

	self.CursorSync()
	self.statusline.Render()
}

func (self *Editor) CursorJumpDown(jump int) {
	self.fileCursor.JumpDown(jump)

	currentRowNumber := self.fileCursor.RowNumberGet() - self.offsetRow
	isOffset := self.rowCount < self.file.RowCount()
	isCrossScrollOff := currentRowNumber >= self.rowCount-self.scrollOff
	isScroll := self.rowCount+self.offsetRow < self.file.RowCount()
	if isOffset &&
		isCrossScrollOff &&
		isScroll {
		self.offsetRow++
		self.Render()
	}

	self.CursorSync()
	self.statusline.Render()
}

func (self *Editor) CursorJumpUp(jump int) {
	self.fileCursor.JumpUp(jump)

	currentRowNumber := self.fileCursor.RowNumberGet() - self.offsetRow
	isOffset := self.rowCount < self.file.RowCount()
	isCrossScrollOff := currentRowNumber <= self.scrollOff
	isScroll := self.offsetRow > 0
	if isOffset &&
		isCrossScrollOff &&
		isScroll {
		self.offsetRow--
		self.Render()
	}

	self.CursorSync()
	self.statusline.Render()
}

func (self *Editor) CursorJumpRight(jump int) {
	self.fileCursor.JumpRight(jump)

	currentColumnNumber := self.fileCursor.ColumnNumberGet() - self.offsetColumn
	isOffset := self.columnCount < self.file.ColumnEnd(self.fileCursor.RowNumberGet())
	isCrossScrollOff := currentColumnNumber >= self.columnCount-self.scrollOff
	isScroll := self.columnCount+self.offsetColumn <= self.file.ColumnEnd(self.fileCursor.RowNumberGet())
	if isOffset &&
		isCrossScrollOff &&
		isScroll {
		self.offsetColumn++
		self.Render()
	}

	self.CursorSync()
	self.statusline.Render()
}

func (self *Editor) CursorSync() {
	self.consoleCommands.CursorSet(
		self.fileCursor.RowNumberGet()-self.offsetRow,
		self.fileCursor.ColumnNumberGet()-self.offsetColumn,
	)
}

func (self *Editor) FileLoad() error {
	if err := self.file.Load(); err != nil {
		return err
	}
	self.consoleCommands.CursorSetTopLeft()
	self.TitleSet(self.file.Name())
	return nil
}

func (self *Editor) FileSave() error {
	err := self.file.Save()
	self.statusline.Render()
	return err
}

func (self *Editor) HeightSet(rowCount int) {
	if self.rowCount == rowCount {
		return
	}
	self.rowCount = rowCount
	self.Render()
}

func (self *Editor) LineAddAbove() {
	self.file.RowAdd(self.fileCursor.RowNumberGet() - 1)
	self.statusline.Render()
	self.Render()
}

func (self *Editor) LineAddBelow() {
	self.file.RowAdd(self.fileCursor.RowNumberGet())
	self.fileCursor.JumpDown(1)
	self.statusline.Render()
	self.Render()
}

func (self *Editor) ModeToEdit() {
	self.editorMode.ToEdit()
	self.statusline.Render()
}

func (self *Editor) ModeToView() {
	self.editorMode.ToView()
	self.statusline.Render()
}

func (self *Editor) RowNumberGet() int {
	return self.fileCursor.RowNumberGet() - self.offsetRow
}

func (self *Editor) RowRender(rowNumber int) {
	self.consoleCommands.CursorHide()

	fileColumn := self.offsetColumn
	fileRowIndex := rowNumber + self.offsetRow - 1
	row := self.file.Rows()[fileRowIndex]
	columnEnd := len(row)
	columnNumber := 1

	for ; columnNumber <= self.columnCount; columnNumber++ {
		fileColumn++
		if fileColumn <= columnEnd {
			self.consoleCommands.RunePrint(rowNumber, columnNumber, row[fileColumn-1])
			continue
		}
		if fileColumn == columnEnd+1 {
			self.consoleCommands.RunePrint(rowNumber, columnNumber, '↲')
			continue
		}
		self.consoleCommands.RunePrint(rowNumber, columnNumber, ' ')
	}
	self.consoleCommands.CursorShow()
}

func (self *Editor) RuneDeleteLeft() {
	self.file.RuneDelete(self.fileCursor.RowNumberGet(), self.fileCursor.ColumnNumberGet())
	self.fileCursor.JumpLeft(1)

	self.RowRender(self.RowNumberGet())
	self.statusline.Render()
}

func (self *Editor) RuneDeleteRight() {
	self.file.RuneDelete(self.fileCursor.RowNumberGet(), self.fileCursor.ColumnNumberGet()+1)

	self.RowRender(self.RowNumberGet())
	self.statusline.Render()
}

func (self *Editor) RuneInsert(char rune) {
	self.file.RuneInsert(self.fileCursor.RowNumberGet(), self.fileCursor.ColumnNumberGet(), char)
	self.fileCursor.JumpRight(1)

	self.RowRender(self.RowNumberGet())
	self.statusline.Render()
}

func (self *Editor) Render() {
	self.consoleCommands.CursorPositionSave()
	rowNumber := 1

	for ; rowNumber <= self.rowCount && rowNumber <= self.file.RowCount(); rowNumber++ {
		self.RowRender(rowNumber)
	}
	for ; rowNumber <= self.rowCount; rowNumber++ {
		self.consoleCommands.RunePrint(rowNumber, 1, '~')
	}
	self.consoleCommands.CursorPositionRestore()
	self.statusline.Render()
}

func (self *Editor) TitleSet(title string) {
	self.consoleCommands.TitleSet("VimGo - " + title)
}

func (self *Editor) WidthSet(columnCount int) {
	if self.columnCount == columnCount {
		return
	}
	self.columnCount = columnCount
	self.Render()
}
