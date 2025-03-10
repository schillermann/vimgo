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
		editor.RuneInsert(keyboard.RuneGet())
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

	switch keyboard.RuneGet() {
	case 'h':
		editor.CursorJumpLeft(1)
	case 'j':
		editor.CursorJumpDown(1)
	case 'k':
		editor.CursorJumpUp(1)
	case 'l':
		editor.CursorJumpRight(1)
	case 'e':
		editor.ModeToEdit()
	case 'o':
		editor.LineAddBelow()
		editor.ModeToEdit()
	case 'O':
		editor.LineAddAbove()
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
	editor := NewEditor(file, fileCursor, editorMode, statusline, consoleCommands, windowRows-statusline.RowCountGet(), windowColumns)
	if err := editor.FileLoad(); err != nil {
		SafeExit(err)
	}
	editor.Render()

	for {
		windowRows, windowColumns, err := consoleWindow.Size()
		if err != nil {
			SafeExit(err)
		}
		statusline.PositionSet(windowRows)
		statusline.WidthSet(windowColumns)

		editor.HeightSet(windowRows - statusline.RowCountGet())
		editor.WidthSet(windowColumns)

		keyboard.Read()
		if editorMode.IsView() {
			modeView(keyboard, editor)
		} else if editorMode.IsEdit() {
			modeEdit(keyboard, editor)
		}
	}
}
