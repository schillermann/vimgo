package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/schillermann/vimgo/console"
)

var consoleWindow = console.Window{}

func KeyPress() (rune, error) {
	input := bufio.NewReader(os.Stdin)
	char, _, err := input.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("error reading key press: %w", err)
	}
	return char, nil
}

func SafeExit(withErr error) {
	consoleCsi.ScreenClear()

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
	if err := consoleWindow.EnableRawMode(); err != nil {
		SafeExit(nil)
	}

	flag.Parse()
	editor := NewEditor(NewFile(flag.Arg(0)))

	if err := editor.FileLoad(); err != nil {
		SafeExit(err)
	}
	if err := editor.ScreenRender(); err != nil {
		SafeExit(err)
	}

	for {
		// key press
		keyPress, err := KeyPress()
		if err != nil {
			SafeExit(err)
		}

		switch keyPress {
		case 'h':
			editor.CursorMoveLeft(1)
		case 'j':
			editor.CursorMoveDown(1)
		case 'k':
			editor.CursorMoveUp(1)
		case 'l':
			editor.CursorMoveRight(1)
		case 'q':
			SafeExit(nil)
		}
	}
}
