//go:build windows
// +build windows

package console

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"golang.org/x/sys/windows"
)

type Window struct {
	settingsBuffer bytes.Buffer
}

type consoleMode struct {
	In  uint32
	Out uint32
}

func (self *Window) DisableRawMode() error {
	var mode consoleMode

	if self.settingsBuffer.Len() == 0 {
		return nil
	}
	if err := gob.NewDecoder(bytes.NewReader(self.settingsBuffer.Bytes())).Decode(&mode); err != nil {
		return fmt.Errorf("error decoding terminal settings: %w", err)
	}

	if err := windows.SetConsoleMode(windows.Stdin, mode.In); err != nil {
		return fmt.Errorf("error setting raw mode for input: %w", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, mode.Out); err != nil {
		return fmt.Errorf("error setting raw mode for output: %w", err)
	}
	return nil
}

func (self *Window) EnableRawMode() error {
	mode := consoleMode{}

	if err := windows.GetConsoleMode(windows.Stdin, &mode.In); err != nil {
		return fmt.Errorf("error getting raw mode for input: %s", err)
	}
	if err := windows.GetConsoleMode(windows.Stdout, &mode.Out); err != nil {
		return fmt.Errorf("error getting raw mode for output: %w", err)
	}

	self.settingsBuffer = bytes.Buffer{}
	if err := gob.NewEncoder(&self.settingsBuffer).Encode(mode); err != nil {
		return fmt.Errorf("error serializing existing console settings: %w", err)
	}

	var inSettings uint32 = windows.ENABLE_EXTENDED_FLAGS | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	var outSettings uint32 = windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_PROCESSED_OUTPUT | windows.DISABLE_NEWLINE_AUTO_RETURN

	if err := windows.SetConsoleMode(windows.Stdin, inSettings); err != nil {
		return fmt.Errorf("error setting raw mode for input: %w", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, outSettings); err != nil {
		return fmt.Errorf("error setting raw mode with output: %s", err)
	}

	return nil
}

func (self *Window) Size() (rows int, columns int, err error) {
	info := windows.ConsoleScreenBufferInfo{}

	if err := windows.GetConsoleScreenBufferInfo(windows.Stdout, &info); err != nil {
		return 0, 0, fmt.Errorf("error fetching screen size: %w", err)
	}

	return int(info.Window.Bottom - info.Window.Top + 1), int(info.Window.Right - info.Window.Left + 1), nil
}
