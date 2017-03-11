package main

import (
	"io"
	"math"
	"os"

	"github.com/daved/simpartsim"
	"github.com/tgreiser/etherdream"
)

type space struct {
	simpartsim.Space
}

func (s *space) run(ps simpartsim.Particles, frames int) <-chan []simpartsim.Coords {
	csc := make(chan []simpartsim.Coords)

	go func() {
		ps.Reset()
		s.Space.Run(ps, frames, csc)

		defer close(csc)
	}()

	return csc
}

func (s *space) toStdout(ps simpartsim.Particles, frames int) error {
	csc := s.run(ps, frames)

	for cs := range csc {
		if err := dumpToStdout(os.Stdout, cs); err != nil {
			return err
		}
	}

	return nil
}

func (s *space) pointStream(ps simpartsim.Particles, frames int) etherdream.PointStream {
	psFn := func(w io.WriteCloser) {
		defer func() {
			_ = w.Close()
		}()

		csc := s.run(ps, frames)
		pts := 0
		fpts := 0
		var last *etherdream.Point
		var err error

		for cs := range csc {
			if pts, last, err = dumpInPointStream(w, cs); err != nil {
				return
			}
			fpts = pts
			times := int(math.Floor(float64(etherdream.FramePoints / pts)))
			for i := 1; i < times; i++ {
				if pts, last, err = dumpInPointStream(w, cs); err != nil {
					return
				}
				fpts += pts
			}
			//log.Printf("Ran %v - %v times for %v\n", pts, times, fpts)
			_ = etherdream.NextFrame(w, fpts, *last)
		}
		w.Close()
	}

	return psFn
}
