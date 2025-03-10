package console

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

type Keyboard struct {
	buffer     []byte
	bufferSize int
}

// Find byte code for key with: showkey -a
var (
	keyBackspace = []byte{127}
	keyDelete    = []byte{27, 91, 51, 126}
	keyEsc       = []byte{27}
	keyPageDown  = []byte{27, 91, 54, 126}
	keyPageUp    = []byte{27, 91, 53, 126}
)

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (self *Keyboard) IsKeyBackspace() bool {
	return bytes.Equal(self.buffer[0:self.bufferSize], keyBackspace)
}

func (self *Keyboard) IsKeyDelete() bool {
	return bytes.Equal(self.buffer[0:self.bufferSize], keyDelete)
}

func (self *Keyboard) IsKeyEsc() bool {
	return bytes.Equal(self.buffer[0:self.bufferSize], keyEsc)
}

func (self *Keyboard) IsRune() bool {
	return unicode.IsPrint(self.RuneGet())
}

func (self *Keyboard) Read() error {
	self.buffer = make([]byte, 4)

	size, err := os.Stdin.Read(self.buffer)
	if err != nil {
		return fmt.Errorf("error reading keyboard input: %w", err)
	}
	self.bufferSize = size
	return nil
}

func (self *Keyboard) RuneGet() rune {
	r, _ := utf8.DecodeRune(self.buffer)
	return r
}
