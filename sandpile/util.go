package sandpile

import (
	"image/color"
)

const (
	Top = iota
	Right
	Left
	Bottom
)

type ColorList []color.RGBA

type point struct {
	x int
	y int
}

type bounds struct {
	minX, maxX int
	minY, maxY int
}

type findOrCreateNodeFunc func(p point) *sandNode
