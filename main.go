package main

import (
	"bufio"
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
	fmt.Fprint(os.Stdout, console.AnsiClearScreen)
	fmt.Fprint(os.Stdout, console.AnsiCursorMoveTopLeft)

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

	// test output
	fmt.Fprint(os.Stdout, console.AnsiClearScreen)
	fmt.Fprintf(os.Stdout, console.AnsiCustorMovePrintRune, 0, 0, 'q')
	fmt.Fprint(os.Stdout, console.AnsiCursorMoveTopLeft)

	keyPress, err := KeyPress()
	if err != nil {
		SafeExit(err)
	}

	switch keyPress {
	case 'q':
		SafeExit(nil)
	}
}
