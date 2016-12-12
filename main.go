package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/daved/simpartsim"
)

func main() {
	stdout := false
	parts := 20
	frames := 200
	opts := simpartsim.SimpleSpaceOptions{
		FrameLen: .1,
		Size:     100.0,
		Gravity:  9.81,
		Drag:     9.0,
	}

	flag.BoolVar(&stdout, "stdout", stdout, "to stdout")
	flag.IntVar(&parts, "parts", parts, "particle count")
	flag.IntVar(&frames, "frames", frames, "frame count")
	flag.Parse()

	spc := space{simpartsim.NewSimpleSpace(opts)}
	ps := simpartsim.NewSimpleParticles(parts, spc.Termination())

	if stdout {
		if err := spc.toStdout(ps, frames); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	stream := spc.pointStream(ps, frames)
	_ = stream // remove when var is used
}
