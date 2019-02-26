package synth

import "github.com/jacoblister/noisefloor/component/synth/processor"

// Patch is a simple minimal example patch
type Patch struct {
	oscillator processor.Oscillator
	envelope   processor.Envelope
	gain       processor.Gain
}

// Start - init patch
func (p *Patch) Start(sampleRate int) {
	p.oscillator.Start(sampleRate)
	p.oscillator.Waveform = processor.Square

	p.envelope.Start(sampleRate)
	p.gain.Start(sampleRate)
}

// Process - produce next sample
func (p *Patch) Process(freq float32, gate float32, trigger float32) (output float32) {
	p.oscillator.Freq = freq
	sample := p.oscillator.Process()
	env := p.envelope.Process(gate, trigger)

	output = p.gain.Process(sample, env)
	return
}