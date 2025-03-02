package main

type EditorMode struct {
	mode string
}

func NewEditorMode() *EditorMode {
	return &EditorMode{
		mode: "view",
	}
}

func (self *EditorMode) IsEdit() bool {
	if self.mode == "edit" {
		return true
	}
	return false
}

func (self *EditorMode) IsView() bool {
	if self.mode == "view" {
		return true
	}
	return false
}
func (self *EditorMode) Mode() string {
	return self.mode
}

func (self *EditorMode) ToEdit() {
	consoleCommands.CursorStyleLine()
	self.mode = "edit"
}

func (self *EditorMode) ToView() {
	consoleCommands.CursorStyleReset()
	self.mode = "view"
}
