package terminal

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"golang.org/x/sys/unix"
)

type TerminalMode struct {
	settings []byte
	termios  *unix.Termios
}

func (self *TerminalMode) Termios() (*unix.Termios, error) {
	if self.termios == nil {
		var termios, err = unix.IoctlGetTermios(unix.Stdin, unix.TCGETS)
		if err != nil {
			return nil, fmt.Errorf("error fetching existing console settings: %w", err)
		}
		self.termios = termios
	}

	return self.termios, nil
}

func (self *TerminalMode) Settings() ([]byte, error) {
	if self.settings == nil {
		var termios, termiosErr = self.Termios()
		if termiosErr != nil {
			return nil, termiosErr
		}

		var buffer = bytes.Buffer{}
		var bufferErr = gob.NewEncoder(&buffer).Encode(termios)
		if bufferErr != nil {
			return nil, fmt.Errorf("error serializing existing console settings: %w", bufferErr)
		}
		self.settings = buffer.Bytes()
	}

	return self.settings, nil
}

func (self *TerminalMode) EnableRawMode() error {
	var termios, termiosErr = self.Termios()
	if termiosErr != nil {
		return termiosErr
	}

	var _, settingsErr = self.Settings()
	if settingsErr != nil {
		return settingsErr
	}

	termios.Lflag = termios.Lflag &^ (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Iflag = termios.Iflag &^ (unix.IXON | unix.ICRNL | unix.BRKINT | unix.INPCK | unix.ISTRIP)
	termios.Oflag = termios.Oflag &^ (unix.OPOST)
	termios.Cflag = termios.Cflag | unix.CS8

	var err = unix.IoctlSetTermios(unix.Stdin, unix.TCSETSF, termios)
	if err != nil {
		return err
	}
	return nil
}

func (self *TerminalMode) DisableRawMode() error {
	var termios unix.Termios
	var settings, settingsErr = self.Settings()
	if settingsErr != nil {
		return settingsErr
	}
	var decodingErr = gob.NewDecoder(bytes.NewReader(settings)).Decode(&termios)
	if decodingErr != nil {
		return fmt.Errorf("error decoding terminal settings: %w", decodingErr)
	}
	var termiosErr = unix.IoctlSetTermios(unix.Stdin, unix.TCSETSF, &termios)
	if termiosErr != nil {
		return fmt.Errorf("error restoring original console settings: %w", termiosErr)
	}
	return nil
}
