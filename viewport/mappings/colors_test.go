package mappings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetForeAndBackgroundFromTileColours(t *testing.T) {
	var inputBytes byte = 0xCC

	var expectedBackgroundBytes byte = 0xC
	var expectedForegroundBytes byte = 0xC

	gotBackgroundBytes, gotForegroundBytes := getCellColoursFromTileColours(inputBytes)

	assert.Equal(t, expectedBackgroundBytes, gotBackgroundBytes, "Background bytes mismatch")
	assert.Equal(t, expectedForegroundBytes, gotForegroundBytes, "Foreground bytes mismatch")
}
