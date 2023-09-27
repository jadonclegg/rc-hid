package control

import (
	"fmt"
	"math"
)

// PotControl is used to control potentiometers, and can map an input resolution of 8-16 bits to an output resolution of 8-16 bits.
type PotControl struct {
	inputResolution  uint16
	outputResolution uint16
	maxInput         uint16
	halfMaxInput     uint16
	Invert           bool
	endpointInfo
}

// NewPotControl initializes a PotControl with the info it needs to work. You need to specify the input and output resolution (8-16 bits)
func NewPotControl(inputResolution, outputResolution uint16) *PotControl {
	control := &PotControl{
		inputResolution:  inputResolution,
		outputResolution: outputResolution,
	}

	control.maxInput = squareUint16(inputResolution) - 1
	control.halfMaxInput = control.maxInput / 2

	control.endpointInfo = *newEndpointInfo(outputResolution)
	fmt.Println("Max output: ", control.maxOutput)
	fmt.Println("Max Input: ", control.maxInput)

	return control
}

// GetOutputValue takes the input value, and converts it to the output value, adjusting for things like resolution difference, output endpoints, and output midpoint
func (c *PotControl) GetOutputValue(input uint16, controlRange Range) uint16 {
	adjustedValue := input
	adjustedValue = uint16(int64(c.maxOutput) * int64(adjustedValue) / int64(c.maxInput))

	if controlRange == LowerHalf {
		adjustedValue = uint16(c.halfMaxOutput - float64(adjustedValue)/float64(2))
	}

	if controlRange == UpperHalf {
		adjustedValue = uint16(c.halfMaxOutput + float64(adjustedValue)/float64(2))
	}

	if c.Invert {
		adjustedValue = c.maxOutput - adjustedValue
	}

	adjustedValue = adjustValueForEndpoints(adjustedValue, c.endpointInfo)

	// // fmt.Println(adjustedValue)

	return adjustedValue
}

func squareUint16(input uint16) uint16 {
	return uint16(math.Pow(2, float64(input)))
}
