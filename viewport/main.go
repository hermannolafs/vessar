package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/kelindar/tile"
	"github.com/logrusorgru/aurora/v3"
	"hermannolafs/vessar/viewport/mappings"
	"log"
	"os"
)

var (
	zeroZeroTile = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x02, 0x51}
	middletile   = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x02, 0x23}
	maxMaxTile   = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x01, 0x51}
	standardTile = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x00, 0x02}

	equalPosTile = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x01, 0x31}
	playerTile   = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x03, 0x14}
)

const (
	// index prefix represents index of data stored in in tile byte array
	indexMapProperties = 4 // 0000 0011 ; 0000 : 00 : 11 ; 0 none 1 collision 2 npc 3 playerc
	indexColor         = 5 // 1111 1111 ; bg 1111 : fg 1111

	playerViewSize = 6 // TODO this should be configurable or hard coded
)

// used for kelidar/tile functions where we do not need to pass a function
func none(_ tile.Point, _ tile.Tile) {}

func main() {
	// the Grid should be read from a grpc request for the map, player pos polled or something cool

	playerView := newPlayerView()

	for {
		playerView.consumeTerminalEvents()
		playerView.printViewToTerminal()
	}
}

func (player *Player) consumeTerminalEvents() {
	event := player.screen.PollEvent()

	switch event := event.(type) {
	case *tcell.EventResize:
		player.screen.Sync()
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyCtrlC:
			player.exit(0)
		case tcell.KeyUp:
			player.MoveNorth()
		}
	}
}

func (player Player) exit(code int) {
	player.screen.Fini()
	os.Exit(code)
}

func (player Player) printViewToTerminal() {
	player.iterateOverPlayerView(func(point tile.Point, t tile.Tile) {
		character := getCharacterForTile(t[indexMapProperties])
		player.screen.SetContent(
			int(point.X),
			int(point.Y),
			character[0],
			character[1:],
			mapGridTileToTcellStyle(t),
		)
	})
	player.screen.Show()
}

func (player *Player) MoveNorth() {
	player.grid.Within(player.position, player.position, func(point tile.Point, t tile.Tile) {
		oldPosition := player.position
		newPosition := point.Move(tile.North)
		player.grid.WriteAt(newPosition.X, newPosition.Y, playerTile)
		player.grid.WriteAt(oldPosition.X, oldPosition.Y, standardTile)
		player.setNewPlayerPosition(newPosition)

	})
}

func getCharacterForTile(mapProperties byte) []rune {
	// This will break if we use any of the upper bits
	// TODO learn how to do bitwise switch cases
	switch mapProperties {
	case mappings.GridNone:
		return []rune{mappings.TerminalRunes[mappings.None]}
	case mappings.GridCollision:
		return []rune{mappings.TerminalRunes[mappings.Collision]}
	case mappings.GridNonPlayer:
		return []rune{mappings.TerminalRunes[mappings.NonPlayer]}
	case mappings.GridPlayer:
		return []rune{mappings.TerminalRunes[mappings.Player]}
	default:
		return []rune{mappings.TerminalRunes[mappings.None]}
	}
}

func mapGridTileToTcellStyle(tile tile.Tile) tcell.Style {
	background, foreground := mappings.GetTerminalColoursFromTileColours(tile[indexColor])
	return tcell.StyleDefault.
		Background(background).
		Foreground(foreground)
}

type Player struct {
	grid     *tile.Grid
	position tile.Point // represents grid size of view
	screen   tcell.Screen
}

func newPlayerView() *Player {
	grid := newGrid(30, 30)
	setReferenceTiles(grid)
	// WIP Setup player character

	playerPosition := getDefaultPlayerPositionSizeAsPoint()
	grid.WriteAt(playerViewSize/2, playerViewSize/2, playerTile)

	terminalScreen := newTcellScreen()

	playerView := &Player{
		grid:     grid,
		position: playerPosition, // Maybe this should be a function instead of private field?
		screen:   terminalScreen,
	}

	return playerView
}

func newTcellScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}

	screen.SetStyle(tcell.StyleDefault)
	return screen
}

func getDefaultPlayerPositionSizeAsPoint() tile.Point {
	return tile.Point{playerViewSize / 2, playerViewSize / 2}
}

// Dumb functiuon??
func (player Player) iterateOverPlayerView(functionToRun func(point tile.Point, tile tile.Tile)) {
	player.grid.Within(
		tile.Point{
			player.position.X - playerViewSize/2,
			player.position.Y - playerViewSize/2,
		},
		tile.Point{
			player.position.X + playerViewSize/2,
			player.position.Y + playerViewSize/2,
		},
		functionToRun,
	)
}

func (player *Player) setNewPlayerPosition(position tile.Point) {
	player.position = position
}

// Calculates the
func calculatePlayerViewRectangleFromPosition(playerPosition tile.Point) tile.Rect {
	return tile.NewRect(
		// Topleft
		playerPosition.X-playerViewSize/2,
		playerPosition.Y-playerViewSize/2,
		// BottomRight
		playerPosition.X+playerViewSize/2,
		playerPosition.Y+playerViewSize/2,
	)
}

// Sets 00, X,Y and X/2,Y/2 as standard tiles
// For developmental purposes, TODO delete/move this
func setReferenceTiles(grid *tile.Grid) {
	grid.Each(func(point tile.Point, _ tile.Tile) {
		grid.WriteAt(point.X, point.Y, standardTile)
	})

	// Marking axis
	for i := int16(0); i < grid.Size.X; i++ {
		grid.WriteAt(i, i, equalPosTile)
	}

	// Mark Corners for reference
	grid.WriteAt(0, 0, zeroZeroTile)
	grid.WriteAt(grid.Size.X-1, grid.Size.Y-1, maxMaxTile)
	// Mark middle
	grid.WriteAt(grid.Size.X/2, grid.Size.Y/2, middletile)
}

// Returns x,y grid
func newGrid(x, y int16) *tile.Grid {

	grid := tile.NewGrid(x, y)
	log.Print("Created Grid of size ", aurora.Green(grid.Size))
	return grid
}
