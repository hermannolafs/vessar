package color

import "github.com/gdamore/tcell"

const (
	ColorMin = 0x00 + iota
	ColorMax
	ColorMiddle

	ColorStandard
	ColorEven
)

var ColorValues = map[byte]tcell.Color{
	ColorMin: 0x00,
	ColorMax: 0x01,
	ColorMiddle: 0x03,

	ColorStandard: 0x04,
	ColorEven: 0x05,
}