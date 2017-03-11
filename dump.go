package main

import (
	"fmt"
	"image/color"
	"io"

	"github.com/daved/simpartsim"
	"github.com/tgreiser/etherdream"
	"github.com/tgreiser/ln/ln"
)

func dumpToStdout(w io.Writer, cs []simpartsim.Coords) error {
	for k := range cs {
		x, y := int(cs[k].X), int(cs[k].Y)
		bs := []byte(fmt.Sprintf("%d,%d\n", x, y))

		if _, err := w.Write(bs); err != nil {
			return err
		}
	}

	return nil
}

func dumpInPointStream(w io.WriteCloser, cs []simpartsim.Coords) (int, *etherdream.Point, error) {
	ct := 0
	var pt *etherdream.Point
	c := color.RGBA{0xff, 0x00, 0x00, 0xff}
	for k := range cs {
		x, y := int(cs[k].X)*(*scale), int(cs[k].Y)*(*scale)
		if pt != nil {
			path := ln.Path{ln.Vector{float64(pt.X), float64(pt.X), 0}, ln.Vector{float64(x), float64(y), 0}}
			etherdream.BlankPath(w, path)
			ct += *etherdream.BlankCount
		}

		pt = etherdream.NewPoint(x, y, c)
		if _, err := w.Write(pt.Encode()); err != nil {
			return ct, pt, err
		}
		ct++
	}

	return ct, pt, nil
}
