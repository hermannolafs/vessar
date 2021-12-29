package mappings

// Map Properties

// subject to change
// tile[4] 0000 0011
const (
	GridNone      byte = 0x0
	GridCollision byte = 0x01
	GridNonPlayer byte = 0x02
	GridPlayer    byte = 0x03
)

const (
	None      byte = 0x0 + iota
	Collision      // TODO there will be more of these
	NonPlayer      // TODO there will be more of these
	Player
)

var TerminalRunes = map[byte]rune{
	None:      '❖',
	Collision: '∆',
	NonPlayer: '♜',
	Player:    '♞',
}
