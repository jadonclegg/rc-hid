package control

import (
	"fmt"
	"math"
)

const DefaultOutputResolution = 12
const DefaultInputResolution = 16

type PotControl struct {
	inputResolution  uint16
	outputResolution uint16
	maxInput         uint16
	halfMaxInput     uint16
	Invert           bool
	endpointInfo
}

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

func (c *PotControl) GetOutputValue(input uint16, controlRange ControlRange) uint16 {
	adjustedValue := input
	adjustedValue = uint16(int64(c.maxOutput) * int64(adjustedValue) / int64(c.maxInput))

	if controlRange == LOWER_HALF {
		adjustedValue = uint16(c.halfMaxOutput - float64(adjustedValue)/float64(2))
	}

	if controlRange == UPPER_HALF {
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
