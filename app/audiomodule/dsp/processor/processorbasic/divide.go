package processorbasic

// Divide - divide x by y
type Divide struct {
	Dummy int
}

// Process - produce next sample
func (d *Divide) Process(x float32, y float32) (Out float32) {
	if y == 0 {
		return 0
	}
	Out = x / y
	return
}
