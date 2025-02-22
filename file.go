package main

import (
	"bufio"
	"fmt"
	"os"
)

type File struct {
	filename string
	rows     [][]rune
}

func NewFile(filename string) *File {
	return &File{
		filename: filename,
	}
}

func (self *File) Filename() string {
	return self.filename
}

func (self *File) NoFile() bool {
	if self.filename == "" {
		return true
	}
	return false
}

func (self *File) Load() error {
	if self.NoFile() {
		return nil
	}

	file, err := os.Open(self.filename)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", self.filename, err)
	}

	rows := [][]rune{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rows = append(rows, []rune(scanner.Text()))
	}
	file.Close()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", self.filename, err)
	}
	self.rows = rows

	return nil
}

func (self *File) Rows() [][]rune {
	return self.rows
}

func (self *File) NumberOfRows() int {
	return len(self.rows)
}
