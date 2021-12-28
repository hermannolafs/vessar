package color

import "github.com/gdamore/tcell/v2"

const (
	maskBackground byte = 0xF0
	maskForeground byte = 0x0F
)

// Auto incremented consts for color mappings
// We only have 16, anything above that will not
// be resolved due to bit masking
const (
	Empty byte = 0x0 + iota
	Purple
	Blue
	Orange
	Yellow
	Green
)

// For input bytes, returns background, foreground as bytes
func getCellColoursFromTileColours(tileColors byte) (byte, byte) {
	background := tileColors & maskBackground >> 4
	foreground := tileColors & maskForeground

	return background, foreground
}

// For input bytes, returns background, foreground as tcell Colour values
func GetTerminalColoursFromTileColours(tileColors byte) (tcell.Color, tcell.Color) {
	background, foreground := getCellColoursFromTileColours(tileColors)

	return TerminalColors[background], TerminalColors[foreground]
}

var TerminalColors = map[byte]tcell.Color{
	// https://colorhunt.co/palette/37066535589af14a16fc9918
	Empty:  tcell.ColorDefault,
	Purple: tcell.NewHexColor(0x370665),
	Blue:   tcell.NewHexColor(0x35589A),
	Orange: tcell.NewHexColor(0xF14A16),
	Yellow: tcell.NewHexColor(0xFC9918),

	Green: tcell.NewHexColor(0x9AE66E),
}
