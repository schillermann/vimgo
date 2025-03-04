package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/schillermann/vimgo/console"
)

var consoleWindow = console.Window{}
var consoleCommands = console.NewCommands()

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

func modeEdit(keyboard *console.Keyboard, editor *Editor) {
	if keyboard.IsRune() {
		editor.RuneInsert(keyboard.GetRune())
		return
	}
	if keyboard.IsKeyDelete() {
		editor.RuneDeleteRight()
		return
	}
	if keyboard.IsKeyEsc() {
		editor.ModeToView()
		return
	}
	if keyboard.IsKeyBackspace() {
		editor.RuneDeleteLeft()
	}
}

func modeView(keyboard *console.Keyboard, editor *Editor) {
	if !keyboard.IsRune() {
		return
	}

	switch keyboard.GetRune() {
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

	windowRows, windowColumns, err := consoleWindow.Size()
	if err != nil {
		SafeExit(err)
	}

	flag.Parse()
	file := NewFile(flag.Arg(0))
	consoleCommands := console.NewCommands()
	fileCursor := NewFileCursor(file, consoleCommands)
	editorMode := NewEditorMode()
	keyboard := console.NewKeyboard()

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

		keyboard.Read()
		if editorMode.IsView() {
			modeView(keyboard, editor)
		} else if editorMode.IsEdit() {
			modeEdit(keyboard, editor)
		}
	}
}
