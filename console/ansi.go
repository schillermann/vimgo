package console

const (
	AnsiEscape              = "\033"
	AnsiCSI                 = AnsiEscape + "[" // control sequence introducer
	AnsiClearScreen         = AnsiCSI + "2J"
	AnsiCursorMoveTopLeft   = AnsiCSI + "H"
	AnsiCustorMovePrintRune = AnsiCSI + "%d;%dH%c"
)
