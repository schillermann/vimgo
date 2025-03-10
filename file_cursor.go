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

func (self *FileCursor) ColumnNumberSet(column int) {
	self.column = column
}

func (self *FileCursor) ColumnNumberGet() int {
	return self.column
}

func (self *FileCursor) JumpDown(jump int) {
	if self.row+jump > self.file.RowCount() {
		self.row = self.file.RowCount()
	} else {
		self.row += jump
	}

	columnEnd := self.file.ColumnEnd(self.row) + 1
	if self.column < columnEnd {
		return
	}
	self.column = columnEnd
}

func (self *FileCursor) JumpLeft(jump int) {
	if self.column-jump < 1 {
		self.column = 1
		return
	}
	self.column -= jump
}

func (self *FileCursor) JumpUp(jump int) {
	if self.row-jump < 1 {
		self.row = 1
	} else {
		self.row -= jump
	}

	columnEnd := self.file.ColumnEnd(self.row) + 1
	if self.column >= columnEnd {
		self.column = columnEnd
	}
}

func (self *FileCursor) JumpRight(jump int) {
	columnEnd := self.file.ColumnEnd(self.row) + 1
	if self.column+jump > columnEnd {
		self.column = columnEnd
		return
	}
	self.column += jump
}

func (self *FileCursor) RowNumberGet() int {
	return self.row
}

func (self *FileCursor) RowSet(row int) {
	self.row = row
}
