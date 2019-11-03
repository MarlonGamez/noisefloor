package processor

import (
	"strconv"
	"strings"

	"github.com/jacoblister/noisefloor/pkg/vdom"
)

const scopeSamples = 100

// Scope - display signal
type Scope struct {
	index      int
	samples    [scopeSamples]float32
	lastSample float32
}

// Start - init Scope
func (s *Scope) Start(sampleRate int) {}

// Process - proccess next sample
func (s *Scope) Process(input float32) {
	if s.index == 0 {
		// wait for zero crossing
		if s.lastSample > 0 || input < 0 {
			s.lastSample = input
			return
		}
	}

	if s.index < scopeSamples {
		s.samples[s.index] = input
		s.index++
	} else {
		s.index = 0
	}
	s.lastSample = input
}

// CustomRenderDimentions get the extended dimentions of the scope
func (s *Scope) CustomRenderDimentions() (width int, height int) {
	return 200, 100
}

// Render - render the scope
func (s *Scope) Render() vdom.Element {
	path := strings.Builder{}
	path.WriteString("M0.5," + strconv.Itoa(int(s.samples[0]*50)+50) + ".5")
	for i := 1; i < scopeSamples; i++ {
		path.WriteString(" L" + strconv.Itoa(i*2) + ".5," + strconv.Itoa(int(s.samples[i]*50)+50) + ".5")
	}

	pathElement := vdom.MakeElement("path",
		"d", path.String(),
		"stroke", "blue",
		"fill", "none",
	)

	return pathElement
}
