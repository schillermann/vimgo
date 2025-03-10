package main

import (
	"fmt"
	"github.com/schillermann/vimgo/console"
	"strings"
)

type Statusline struct {
	file              *File
	fileCursor        *FileCursor
	editorMode        *EditorMode
	consoleCommands   *console.Commands
	positionRowNumber int
	columnCount       int
}

func NewStatusline(
	file *File,
	fileCursor *FileCursor,
	editorMode *EditorMode,
	consoleCommands *console.Commands,
	positionRowNumber int,
	columnCount int,
) *Statusline {
	return &Statusline{
		file:              file,
		fileCursor:        fileCursor,
		editorMode:        editorMode,
		consoleCommands:   consoleCommands,
		positionRowNumber: positionRowNumber,
		columnCount:       columnCount,
	}
}

func (self *Statusline) PositionSet(rowNumber int) {
	self.positionRowNumber = rowNumber
}

func (self *Statusline) Render() {
	statuslineMiddleColumns := 1
	statusline := self.String(statuslineMiddleColumns, 0)
	statuslineColumns := len(statusline)
	if statuslineColumns < self.columnCount {
		statusline = self.String(self.columnCount+statuslineMiddleColumns-statuslineColumns, 0)
	} else if statuslineColumns > self.columnCount {
		statusline = self.String(statuslineMiddleColumns, statuslineColumns-self.columnCount)
	}

	self.consoleCommands.CursorPositionSave()
	self.consoleCommands.CursorHide()
	self.consoleCommands.ColorInverse()
	for index, char := range statusline {
		self.consoleCommands.RunePrint(self.positionRowNumber, index+1, char)
	}
	self.consoleCommands.CursorShow()
	self.consoleCommands.CursorPositionRestore()
	self.consoleCommands.ResetFormatting()
}

func (self *Statusline) RowCountGet() int {
	return 1
}

func (self *Statusline) String(middleColumns int, columnsCut int) string {
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
		strings.ToUpper(self.editorMode.Mode()),
		filename,
		fileModified,
		strings.Repeat(" ", middleColumns),
		self.fileCursor.RowNumberGet(),
		self.file.RowCount(),
		self.fileCursor.ColumnNumberGet(),
	)
}

func (self *Statusline) WidthSet(columnCount int) {
	self.columnCount = columnCount
}
