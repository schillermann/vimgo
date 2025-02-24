package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type File struct {
	filename string
	modified bool
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

func (self *File) Insert(row int, column int, char rune) {
	self.rows[row-1] = slices.Insert(self.rows[row-1], column-1, char)
	self.modified = true
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

func (self *File) Modified() bool {
	return self.modified
}

func (self *File) NoFile() bool {
	if self.filename == "" {
		return true
	}
	return false
}

func (self *File) NumberOfRows() int {
	return len(self.rows)
}

func (self *File) Rows() [][]rune {
	return self.rows
}

func (self *File) Save() error {
	file, err := os.Create(self.filename)
	if err != nil {
		file.Close()
		return fmt.Errorf("error creating file: %s: %w", self.filename, err)
	}

	for _, row := range self.rows {
		_, err = file.WriteString(string(row) + "\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("error writing to file %s: %w", self.filename, err)
		}
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("error closing written file: %s: %w", self.filename, err)
	}
	self.modified = false
	return nil
}
