package main

import (
	"flag"
	"image/color"
	"log"
	"os"
	"runtime/pprof"

	"github.com/kleary/sandpiles/sandpile"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	config := &sandpile.Config{
		NumGrains: 10000,
		MaxGrains: 3,
		Colors: sandpile.ColorList{
			color.RGBA{0, 0, 0, 255},
			color.RGBA{61, 144, 34, 255},
			color.RGBA{164, 169, 40, 255},
			color.RGBA{31, 75, 109, 255},
		},
		NumWorkers: 16,
		FileName:   "out3.png",
	}

	sandpile := sandpile.NewSandpile(config)
	sandpile.Process()
}
