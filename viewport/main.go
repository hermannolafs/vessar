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
	player       = tile.Tile{0x00, 0x00, 0x00, 0x00, 0x04, 0x54}
)

const (
	// index prefix represents index of data stored in in tile byte array
	indexMapProperties = 4 // 0000 0011 ; 0000 : 00 : 11 ; 0 none 1 collision 2 npc 3 playerc
	indexColor         = 5 // 1111 1111 ; bg 1111 : fg 1111

	defaultPlayerViewSize = 6
)

func main() {
	// the Grid should be read from a grpc request for the map, player pos polled or something cool
	grid := newGrid(30, 30)
	log.Print("Created Grid of size ", aurora.Green(grid.Size))

	playerView := setupPlayerView(grid)

	for {
		playerView.consumeTerminalEvents()
		playerView.printViewToTerminal()
	}
}

func (playerView PlayerView) consumeTerminalEvents() {
	event := playerView.screen.PollEvent()

	switch event := event.(type) {
	case *tcell.EventResize:
		playerView.screen.Sync()
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyCtrlC:
			playerView.exit(0)
		case tcell.KeyUp:

		}
	}
}

func (playerView PlayerView) exit(code int) {
	playerView.screen.Fini()
	os.Exit(code)
}

func (playerView PlayerView) printViewToTerminal() {
	playerView.view.Each(func(point tile.Point, t tile.Tile) {
		playerView.screen.SetContent(
			int(point.X),
			int(point.Y),
			getCharacterForTile(t[indexMapProperties]),
			nil,
			mapGridTileToTcellStyle(t),
		)
	})
	playerView.screen.Show()
}

func getCharacterForTile(mapProperties byte) rune {
	// This will break if we use any of the upper bits
	// TODO learn how to do bitwise switch cases
	switch mapProperties {
	case mappings.GridNone:
		return mappings.TerminalRunes[mappings.None]
	case mappings.GridCollision:
		return mappings.TerminalRunes[mappings.Collision]
	case mappings.GridNonPlayer:
		return mappings.TerminalRunes[mappings.NonPlayer]
	case mappings.GridPlayer:
		return mappings.TerminalRunes[mappings.Player]
	default:
		return mappings.TerminalRunes[mappings.None]
	}
}

func mapGridTileToTcellStyle(tile tile.Tile) tcell.Style {
	background, foreground := mappings.GetTerminalColoursFromTileColours(tile[indexColor])
	return tcell.StyleDefault.
		Background(background).
		Foreground(foreground)
}

func setupPlayerView(grid *tile.Grid) PlayerView {
	setReferenceTiles(grid)
	return newPlayerView(grid)
}

type PlayerView struct {
	view   *tile.View
	size   tile.Point // represents grid size of view
	screen tcell.Screen
}

func newPlayerView(grid *tile.Grid) PlayerView {
	sizePos := getDefaultPlayerViewSizeAsPoint()

	// WIP Setup player character
	grid.WriteAt(sizePos.X/2, sizePos.Y/2, player)

	tileView := newPlayerViewFromGrid(grid, sizePos)
	terminalScreen := newTcellScreen()

	playerView := PlayerView{
		size:   sizePos, // Maybe this should be a function instead of private field?
		view:   tileView,
		screen: terminalScreen,
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

func getDefaultPlayerViewSizeAsPoint() tile.Point {
	return tile.Point{defaultPlayerViewSize, defaultPlayerViewSize}
}

func newPlayerViewFromGrid(grid *tile.Grid, sizePos tile.Point) *tile.View {
	tileView := grid.View(
		rectFromTwoPositions(
			tile.Point{0, 0},
			sizePos,
		),
		func(p tile.Point, tile tile.Tile) {},
	)

	return tileView
}

// (x1,y1), (x2,y2) --> rect(x1,y1,x2,y2)
// Wraps tile.New:w
func rectFromTwoPositions(lowPosition, highPosition tile.Point) tile.Rect {
	return tile.NewRect(
		lowPosition.X, lowPosition.Y,
		highPosition.X, highPosition.Y,
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
	return tile.NewGrid(x, y)
}
