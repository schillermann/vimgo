package main

import (
	"fmt"
	"github.com/schillermann/vimgo/console"
	"strings"
)

type Statusline struct {
	file            *File
	fileCursor      *FileCursor
	editorMode      *EditorMode
	consoleCommands *console.Commands
	rowsPosition    int
	columnsWidth    int
}

func NewStatusline(
	file *File,
	fileCursor *FileCursor,
	editorMode *EditorMode,
	consoleCommands *console.Commands,
	rowsPosition int,
	columnsWidth int,
) *Statusline {
	return &Statusline{
		file:            file,
		fileCursor:      fileCursor,
		editorMode:      editorMode,
		consoleCommands: consoleCommands,
		rowsPosition:    rowsPosition,
		columnsWidth:    columnsWidth,
	}
}

func (self *Statusline) GetRowsHeight() int {
	return 1
}

func (self *Statusline) Render() {
	self.consoleCommands.ColorInverse()

	statuslineMiddleColumns := 1
	statusline := self.String(statuslineMiddleColumns, 0)
	statuslineColumns := len(statusline)
	if statuslineColumns < self.columnsWidth {
		statusline = self.String(self.columnsWidth+statuslineMiddleColumns-statuslineColumns, 0)
	} else if statuslineColumns > self.columnsWidth {
		statusline = self.String(statuslineMiddleColumns, statuslineColumns-self.columnsWidth)
	}

	self.consoleCommands.CursorHide()
	for index, char := range statusline {
		self.consoleCommands.RunePrint(self.rowsPosition, index+1, char)
	}
	self.consoleCommands.CursorShow()
	self.consoleCommands.ResetFormatting()
	self.consoleCommands.CursorMoveTo(self.fileCursor.GetRow(), self.fileCursor.GetColumn())
}

func (self *Statusline) SetColumnsWidth(columns int) {
	self.columnsWidth = columns
}

func (self *Statusline) SetRowsPosition(row int) {
	self.rowsPosition = row
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
		self.fileCursor.GetRow(),
		self.file.NumberOfRows(),
		self.fileCursor.GetColumn(),
	)
}
