//go:build unix
// +build unix

package console

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"golang.org/x/sys/unix"
)

type Window struct {
	settingsBuffer bytes.Buffer
}

func (self *Window) DisableRawMode() error {
	var termios unix.Termios

	if self.settingsBuffer.Len() == 0 {
		return nil
	}
	if err := gob.NewDecoder(bytes.NewReader(self.settingsBuffer.Bytes())).Decode(&termios); err != nil {
		return fmt.Errorf("error decoding terminal settings: %w", err)
	}
	if err := unix.IoctlSetTermios(unix.Stdin, unix.TCSETSF, &termios); err != nil {
		return fmt.Errorf("error restoring original console settings: %w", err)
	}
	return nil
}

func (self *Window) EnableRawMode() error {
	termios, err := unix.IoctlGetTermios(unix.Stdin, unix.TCGETS)
	if err != nil {
		return fmt.Errorf("error fetching existing console settings: %w", err)
	}

	self.settingsBuffer = bytes.Buffer{}
	if err := gob.NewEncoder(&self.settingsBuffer).Encode(termios); err != nil {
		return fmt.Errorf("error serializing existing console settings: %w", err)
	}

	termios.Lflag = termios.Lflag &^ (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Iflag = termios.Iflag &^ (unix.IXON | unix.ICRNL | unix.BRKINT | unix.INPCK | unix.ISTRIP)
	termios.Oflag = termios.Oflag &^ (unix.OPOST)
	termios.Cflag = termios.Cflag | unix.CS8

	if err := unix.IoctlSetTermios(unix.Stdin, unix.TCSETSF, termios); err != nil {
		return err
	}
	return nil
}

func (self *Window) Size() (columns int, rows int, err error) {
	winsize, err := unix.IoctlGetWinsize(unix.Stdout, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, fmt.Errorf("error fetching window size: %w", err)
	}

	return int(winsize.Col), int(winsize.Row), nil
}
