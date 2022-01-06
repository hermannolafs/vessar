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

	viewSize = 9  // TODO this should be configurable or hard coded
	gridSize = 30 // gridSize X gridSize Always assuming grids are complete rectangle
)

// used for kelidar/tile functions where we do not need to pass a function
func none(_ tile.Point, _ tile.Tile) {}

func main() {
	// the Grid should be read from a grpc request for the map, player pos polled or something cool
	playerView := newPlayerView()

	for {
		playerView.printViewToTerminal()
		playerView.consumeTerminalEvents()
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
			player.MoveSouth()
		case tcell.KeyDown:
			player.MoveNorth()
		case tcell.KeyLeft:
			player.MoveWest()
		case tcell.KeyRight:
			player.MoveEast()
		}
	}
	player.screen.Clear()
}

func (player Player) exit(code int) {
	player.screen.Fini()
	os.Exit(code)
}

func (player Player) printViewToTerminal() {
	player.updatePlayerView()
	player.screen.Show()
}

func (player Player) setPointToTile(point tile.Point, t tile.Tile) {
	// Get character and its properties for tile
	character := getCharacterForTile(t[indexMapProperties])

	// Set terminal point to character found in tile properties
	player.screen.SetContent(
		int(point.X),
		int(point.Y),
		character[0], character[1:],
		mapGridTileToTcellStyle(t),
	)
}

func (player *Player) MoveNorth() { player.MoveInDirection(tile.North) }
func (player *Player) MoveSouth() { player.MoveInDirection(tile.South) }
func (player *Player) MoveEast()  { player.MoveInDirection(tile.East) }
func (player *Player) MoveWest()  { player.MoveInDirection(tile.West) }

func (player *Player) MoveInDirection(direction tile.Direction) {
	player.grid.Within(player.position, player.position, func(point tile.Point, t tile.Tile) {
		// WIP replacement should work differently, this way everywhere the player has been
		// gets overwritten by standard tile.
		oldPosition := player.position
		newPosition := point.Move(direction)
		if isPointOutOfBounds(newPosition) {
			return
		}

		player.grid.WriteAt(newPosition.X, newPosition.Y, playerTile)
		player.grid.WriteAt(oldPosition.X, oldPosition.Y, standardTile)
		player.setNewPlayerPosition(newPosition)
	})
}

func isPointOutOfBounds(position tile.Point) bool {
	switch {
	case position.X < 0:
		return true
	case gridSize <= position.X:
		return true
	case position.Y < 0:
		return true
	case gridSize <= position.Y:
		return true
	}
	return false
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

func getDefaultPlayerPositionSizeAsPoint() tile.Point {
	return tile.Point{X: viewSize / 2, Y: viewSize / 2}
}

func newPlayerView() *Player {
	grid := newGrid(gridSize, gridSize)
	setReferenceTiles(grid)

	playerPosition := getDefaultPlayerPositionSizeAsPoint()
	grid.WriteAt(playerPosition.X, playerPosition.Y, playerTile)

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

// This function is a bit doomed to be long due to the offset check
// needs to exist outside the iterator function passed to the grid.Within(),
// but I don't like global vars for this sort of stuff.
// Maybe when I refactor the top level struct to improve cohesion
// we can put the offset check in a private member on some sub-struct 🤷‍♀️
func (player Player) updatePlayerView() {
	topLeft, bottomRight := player.extractPlayerViewCorners()

	// Represents the point values for the first point.
	// the way the tile lib iterates over the square it will always begin with
	// the "0,0" coordinates for said square, so we store those coordinates
	// and use as offset to where in the terminal to display that exact point.
	// This means our character will always be in the center of the screen.
	// TODO give sense of movement !!?
	offset := tile.Point{}
	offsetFound := false

	player.grid.Within(
		topLeft,
		bottomRight,
		func(point tile.Point, tileAtPoint tile.Tile) {
			if offsetFound == false {
				offset = point
				offsetFound = true
			}

			player.setTerminalToTileAtPoint(
				tileAtPoint,
				tile.Point{
					X: point.X - offset.X,
					Y: point.Y - offset.Y,
				},
			)
		},
	)
}

func (player Player) setTerminalToTileAtPoint(t tile.Tile, point tile.Point) {
	// Get character and its properties for tile
	character := getCharacterForTile(t[indexMapProperties])

	// Set terminal at point x,y with character
	player.screen.SetContent(
		int(point.X), int(point.Y),
		character[0], character[1:],
		mapGridTileToTcellStyle(t),
	)
}

// Iterates over grid representing the players view, of size gridSize X gridSize
func (player Player) iterateOverPlayerView(functionToRun func(point tile.Point, tile tile.Tile)) {
	topLeft, bottomRight := player.extractPlayerViewCorners()
	player.grid.Within(
		topLeft,
		bottomRight,
		functionToRun,
	)
}

// TODO make this function shorter this is a mess
// returns topleft, bottomright
func (player Player) extractPlayerViewCorners() (topLeft, bottomRight tile.Point) {
	// Redundant return for readability
	return calculateViewTopLeftOffset(player.position),
		calculateViewBottomRightOffset(player.position)
}

func calculateViewBottomRightOffset(playerPosition tile.Point) tile.Point {
	// No index issue here so far, just for safety
	bottomRight := tile.Point{
		X: playerPosition.X + viewSize/2,
		Y: playerPosition.Y + viewSize/2,
	}
	// Prevent trying to render the grid on max(x|y)
	if gridSize < playerPosition.X+viewSize/2 {
		bottomRight.X = gridSize
	}
	if gridSize < playerPosition.Y+viewSize/2 {
		bottomRight.Y = gridSize
	}

	return bottomRight
}

func calculateViewTopLeftOffset(playerPosition tile.Point) tile.Point {
	// Going below zero can cause weird overflow, something with tiles lib?
	topLeft := tile.Point{
		X: playerPosition.X - viewSize/2,
		Y: playerPosition.Y - viewSize/2,
	}
	// Prevent trying to render the grid on min(x|y)
	if playerPosition.X < viewSize/2 {
		topLeft.X = 0
	}
	if playerPosition.Y < viewSize/2 {
		topLeft.Y = 0
	}

	return topLeft
}

func (player *Player) setNewPlayerPosition(position tile.Point) {
	// TODO maybe to some move verification here?
	player.position = position
}

// Sets 00, X=Y and X/2,Y/2 as nonstandard tiles
// For developmental purposes, TODO delete/move this
func setReferenceTiles(grid *tile.Grid) {
	grid.Each(func(point tile.Point, _ tile.Tile) {
		grid.WriteAt(point.X, point.Y, standardTile)
	})

	// Marking axis
	for i := int16(0); i < grid.Size.X; i++ {
		grid.WriteAt(i, i, equalPosTile)
	}

	grid.WriteAt(0, 0, zeroZeroTile)
	grid.WriteAt(grid.Size.X-1, grid.Size.Y-1, maxMaxTile)
	grid.WriteAt(grid.Size.X/2, grid.Size.Y/2, middletile)
}

// Returns x,y grid
func newGrid(x, y int16) *tile.Grid {
	if x%3 != 0 || y%3 != 0 {
		panic("grid size needs to be multiple of 3")
	}

	grid := tile.NewGrid(x, y)
	log.Print("Created Grid of size ", aurora.Green(grid.Size))
	return grid
}
