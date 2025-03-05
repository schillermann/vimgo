package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type File struct {
	name     string
	modified bool
	rows     [][]rune
}

func NewFile(name string) *File {
	return &File{
		name: name,
	}
}

func (self *File) ColumnEnd(row int) int {
	return len(self.rows[row-1])
}

func (self *File) Load() error {
	if self.NoFile() {
		return nil
	}

	file, err := os.Open(self.name)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", self.name, err)
	}

	rows := [][]rune{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rows = append(rows, []rune(scanner.Text()))
	}
	file.Close()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", self.name, err)
	}
	self.rows = rows

	return nil
}

func (self *File) Modified() bool {
	return self.modified
}

func (self *File) Name() string {
	return self.name
}

func (self *File) NoFile() bool {
	if self.name == "" {
		return true
	}
	return false
}

func (self *File) NumberOfRows() int {
	return len(self.rows)
}

func (self *File) RowAdd(row int) {
	self.rows = append(self.rows[:row], append([][]rune{{}}, self.rows[row:]...)...)
	self.modified = true
}

func (self *File) Rows() [][]rune {
	return self.rows
}

func (self *File) RuneDelete(row int, column int) {
	self.rows[row-1] = slices.Delete(self.rows[row-1], column-2, column-1)
	self.modified = true
}

func (self *File) RuneInsert(row int, column int, char rune) {
	self.rows[row-1] = slices.Insert(self.rows[row-1], column-1, char)
	self.modified = true
}

func (self *File) Save() error {
	file, err := os.Create(self.name)
	if err != nil {
		file.Close()
		return fmt.Errorf("error creating file: %s: %w", self.name, err)
	}

	for _, row := range self.rows {
		_, err = file.WriteString(string(row) + "\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("error writing to file %s: %w", self.name, err)
		}
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("error closing written file: %s: %w", self.name, err)
	}
	self.modified = false
	return nil
}
