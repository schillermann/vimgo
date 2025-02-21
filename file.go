package main

import (
	"bufio"
	"fmt"
	"os"
)

type File struct {
	Filename string
	rows     [][]rune
}

func (self *File) Read() error {
	file, err := os.Open(self.Filename)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", self.Filename, err)
	}

	rows := [][]rune{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rows = append(rows, []rune(scanner.Text()))
	}
	file.Close()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", self.Filename, err)
	}
	self.rows = rows

	return nil
}

func (self *File) Rows() [][]rune {
	return self.rows
}
