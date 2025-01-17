package processor

//Connector specifies a connection between two Processors
type Connector struct {
	FromProcessor Processor
	FromPort      int
	ToProcessor   Processor
	ToPort        int

	Value   float32   // current sample value
	samples []float32 // current samples
}

//Processor is the getter for the Connector Processor
func (c *Connector) Processor(isInput bool) Processor {
	if isInput {
		return c.ToProcessor
	}
	return c.FromProcessor
}

//Port is the getter for the Connector Port
func (c *Connector) Port(isInput bool) int {
	if isInput {
		return c.ToPort
	}
	return c.FromPort
}

//SetProcessor is the setter for the Connector Processor
func (c *Connector) SetProcessor(isInput bool, processor Processor) {
	if isInput {
		c.ToProcessor = processor
		return
	}
	c.FromProcessor = processor
}

//SetPort is the setter for the Connector Port
func (c *Connector) SetPort(isInput bool, port int) {
	if isInput {
		c.ToPort = port
		return
	}
	c.FromPort = port
}

//Samples is the getter for audio samples
func (c *Connector) Samples() []float32 {
	return c.samples
}

//SetSamples is the setter for audio samples
func (c *Connector) SetSamples(samples []float32) {
	c.samples = samples
}
