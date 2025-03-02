package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/schillermann/vimgo/console"
)

var consoleWindow = console.Window{}
var consoleCommands = console.NewCommands()

func KeyPress() (rune, error) {
	input := bufio.NewReader(os.Stdin)
	char, _, err := input.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("error reading key press: %w", err)
	}
	return char, nil
}

func SafeExit(withErr error) {
	consoleCommands.Reset()

	if err := consoleWindow.DisableRawMode(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: disabling raw mode: %s\r\n", err)
	}

	if withErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\r\n", withErr)
		os.Exit(1)
	}
	os.Exit(0)
}

func modeEdit(editor *Editor, keyPress rune) error {
	switch keyPress {
	case 27: // ESC key
		editor.ModeToView()
	case 127: // Backspace
		editor.RuneDelete()
	default:
		editor.RuneInsert(keyPress)
	}

	return nil
}

func modeView(editor *Editor, keyPress rune) {
	switch keyPress {
	case 'h':
		editor.CursorMoveLeft(1)
	case 'j':
		editor.CursorMoveDown(1)
	case 'k':
		editor.CursorMoveUp(1)
	case 'l':
		editor.CursorMoveRight(1)
	case 'e':
		editor.ModeToEdit()
	case 's':
		if err := editor.FileSave(); err != nil {
			SafeExit(err)
		}
	case 'q':
		SafeExit(nil)
	}
}

func main() {
	if err := consoleWindow.EnableRawMode(); err != nil {
		SafeExit(nil)
	}

	flag.Parse()
	file := NewFile(flag.Arg(0))
	consoleCommands := console.NewCommands()
	fileCursor := NewFileCursor(file, consoleCommands)
	editorMode := NewEditorMode()

	windowRows, windowColumns, err := consoleWindow.Size()
	if err != nil {
		SafeExit(err)
	}
	statusline := NewStatusline(file, fileCursor, editorMode, consoleCommands, windowRows, windowColumns)
	editor := NewEditor(file, fileCursor, editorMode, statusline, consoleCommands, windowRows-statusline.GetRowsHeight(), windowColumns)
	if err := editor.FileLoad(); err != nil {
		SafeExit(err)
	}
	editor.ScreenRender()

	for {
		windowRows, windowColumns, err := consoleWindow.Size()
		if err != nil {
			SafeExit(err)
		}
		statusline.SetRowsPosition(windowRows)
		statusline.SetColumnsWidth(windowColumns)

		editor.SetRowsHeight(windowRows - statusline.GetRowsHeight())
		editor.SetColumnsWidth(windowColumns)

		// key press
		keyPress, err := KeyPress()
		if err != nil {
			SafeExit(err)
		}

		if editorMode.IsView() {
			modeView(editor, keyPress)
		} else if editorMode.IsEdit() {
			if err := modeEdit(editor, keyPress); err != nil {
				SafeExit(err)
			}
		}
	}
}
