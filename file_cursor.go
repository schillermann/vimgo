package main

import (
	"github.com/schillermann/vimgo/console"
)

type FileCursor struct {
	file            *File
	consoleCommands *console.Commands
	row             int
	column          int
}

func NewFileCursor(file *File, consoleCommands *console.Commands) *FileCursor {
	return &FileCursor{
		file:            file,
		consoleCommands: consoleCommands,
		row:             1,
		column:          1,
	}
}

func (self *FileCursor) GetColumn() int {
	return self.column
}

func (self *FileCursor) GetRow() int {
	return self.row
}

func (self *FileCursor) MoveDown(jump int) {
	if self.row+jump > len(self.file.Rows()) {
		return
	}
	self.row += jump

	columnEnd := self.file.ColumnEnd(self.row) + 1
	if self.column > columnEnd {
		self.column = columnEnd
	}
}

func (self *FileCursor) MoveLeft(jump int) {
	if self.column-jump < 1 {
		return
	}
	self.column -= jump
}

func (self *FileCursor) MoveUp(jump int) {
	if self.row-jump < 1 {
		return
	}
	self.row -= jump

	columnEnd := self.file.ColumnEnd(self.row) + 1
	if self.column > columnEnd {
		self.column = columnEnd
	}
}

func (self *FileCursor) MoveRight(jump int) {
	if self.column+jump > self.file.ColumnEnd(self.row)+1 {
		return
	}
	self.column += jump
}

func (self *FileCursor) SetColumn(column int) {
	self.column = column
}

func (self *FileCursor) SetRow(row int) {
	self.row = row
}
