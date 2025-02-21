package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/schillermann/vimgo/console"
)

var consoleWindow = console.Window{}
var consoleCsi = console.Csi{Writer: os.Stdout}

func KeyPress() (rune, error) {
	input := bufio.NewReader(os.Stdin)
	char, _, err := input.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("error reading key press: %w", err)
	}
	return char, nil
}

func SafeExit(withErr error) {
	consoleCsi.ClearScreen()

	if err := consoleWindow.DisableRawMode(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: disabling raw mode: %s\r\n", err)
	}

	if withErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\r\n", withErr)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	flag.Parse()

	if err := consoleWindow.EnableRawMode(); err != nil {
		SafeExit(nil)
	}

	// file output
	file := File{Filename: flag.Arg(0)}
	if err := file.Read(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	consoleCsi.ClearScreen()

	rows, columns, err := consoleWindow.Size()
	if err != nil {
		SafeExit(err)
	}

	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		if rowIndex >= len(file.Rows()) {
			consoleCsi.PrintRune(rowIndex+1, 1, '~')
			continue
		}

		for columnIndex, char := range file.Rows()[rowIndex] {
			if columnIndex >= columns {
				break
			}
			consoleCsi.PrintRune(rowIndex+1, columnIndex+1, char)
		}
	}

	consoleCsi.MoveCursorLeftCorner()

	for {
		// key press
		keyPress, err := KeyPress()
		if err != nil {
			SafeExit(err)
		}

		switch keyPress {
		case 'h':
			consoleCsi.MoveCursorLeft(1)
		case 'j':
			consoleCsi.MoveCursorDown(1)
		case 'k':
			consoleCsi.MoveCursorUp(1)
		case 'l':
			consoleCsi.MoveCursorRight(1)
		case 'q':
			SafeExit(nil)
		}
	}
}
