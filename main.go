package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/schillermann/vimgo/terminal"
)

var terminalMode = terminal.TerminalMode{}
var reader = bufio.NewReader(os.Stdin)

func keyPress() (rune, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("error reading key press: %w", err)
	}
	return char, nil
}

func main() {
	var err = terminalMode.EnableRawMode()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for {
		var keyPress, keyPressErr = keyPress()
		if keyPressErr != nil {
			fmt.Fprintln(os.Stderr, keyPressErr)
			os.Exit(1)
		}

		switch keyPress {
		case 'q':
			var disableRawModeErr = terminalMode.DisableRawMode()
			if disableRawModeErr != nil {
				fmt.Fprintf(os.Stderr, "error disabling terminal raw mode: %v\n", disableRawModeErr)
				os.Exit(1)
			}
			os.Exit(0)
		default:
			fmt.Println(keyPress)
		}
	}
}
