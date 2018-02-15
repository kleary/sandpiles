package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type point struct {
	x int
	y int
}

var (
	maxX int
	minX int
	maxY int
	minY int

	maxGrains = 3

	sand = make(map[point]int)

	notDone = true
	done    = false

	colors = []color.RGBA{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	}
)

func main() {

	sand[point{}] = 200000

	for notDone {
		//topple
		notDone = false
		for p, grains := range sand {
			if grains > maxGrains {
				up := p.y + 1
				right := p.x + 1
				down := p.y - 1
				left := p.x - 1

				if up > maxY {
					maxY = up
				}
				if right > maxX {
					maxX = right
				}
				if down < minY {
					minY = down
				}
				if left < minX {
					minX = left
				}

				val := sand[p] - 4
				sand[p] = val
				notDone = (val > maxGrains) || notDone

				val = sand[point{p.x, up}] + 1
				sand[point{p.x, up}] = val
				notDone = (val > maxGrains) || notDone

				val = sand[point{right, p.y}] + 1
				sand[point{right, p.y}] = val
				notDone = (val > maxGrains) || notDone

				val = sand[point{p.x, down}] + 1
				sand[point{p.x, down}] = val
				notDone = (val > maxGrains) || notDone

				val = sand[point{left, p.y}] + 1
				sand[point{left, p.y}] = val
				notDone = (val > maxGrains) || notDone
			}
		}

	}

	//fmt.Println(sand, minX, maxX, minY, maxY)
	offsetX := -minX
	offsetY := -minY
	img := image.NewRGBA(image.Rect(0, 0, maxX+offsetX, maxY+offsetY))

	for x := 0; x < maxX+offsetX; x++ {
		for y := 0; y < maxY+offsetY; y++ {
			img.Set(x, y, colors[sand[point{x - offsetX, y - offsetY}]])
		}
	}

	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)

}
