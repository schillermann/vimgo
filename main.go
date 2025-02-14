package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/schillermann/vimgo/terminal"
)

var terminalWindow = terminal.TerminalWindow{}
var reader = bufio.NewReader(os.Stdin)

/*** output ***/
func editorRefreshScreen() {
	// clear screen
	fmt.Fprint(os.Stdout, "\x1b[2J")
	fmt.Fprint(os.Stdout, "\x1b[H")

	editorDrawRows()

	// reposition cursor
	fmt.Fprint(os.Stdout, "\x1b[H")
}

func editorDrawRows() {
	var _, rows, err = terminalWindow.NumberOfColumnsAndRows()
	if err != nil {
		exit(err)
	}
	for j := 0; j < rows; j++ {
		fmt.Fprint(os.Stdout, "~")
		if j < rows-1 {
			fmt.Fprint(os.Stdout, "\r\n")
		}
	}
}

/*** Input ***/
func keyPress() (rune, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("error reading key press: %w", err)
	}
	return char, nil
}

func exit(withErr error) {
	// clear screen and
	fmt.Fprint(os.Stdout, "\x1b[2J")
	// reposition cursor to home
	fmt.Fprint(os.Stdout, "\x1b[H")

	var disableRawModeErr = terminalWindow.DisableRawMode()
	if disableRawModeErr != nil {
		fmt.Fprintf(os.Stderr, "error disabling terminal raw mode: %v\n", disableRawModeErr)
		os.Exit(1)
	}

	if withErr == nil {
		os.Exit(0)
	}

	fmt.Fprint(os.Stderr, withErr)
	os.Exit(1)
}

func main() {
	var enableRawModeErr = terminalWindow.EnableRawMode()
	if enableRawModeErr != nil {
		exit(enableRawModeErr)
	}

	for {
		editorRefreshScreen()

		var keyPress, keyPressErr = keyPress()
		if keyPressErr != nil {
			exit(keyPressErr)
		}

		switch keyPress {
		case 'q':
			exit(nil)
		default:
			fmt.Println(keyPress)
		}
	}
}
