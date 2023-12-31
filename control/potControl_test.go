package control

import "testing"

func TestCreatePotControl(t *testing.T) {
	control := NewPotControl(16, 12)
	if control.maxInput != 65535 {
		t.Errorf("when input resolution is 16 bits, max input should be 65535")
	}

	if control.maxOutput != 4095 {
		t.Errorf("when output resolution is 12 bits, max output should be 4095")
	}

	if control.lowerEndpoint != 0 {
		t.Errorf("lower endpoint should be set to 0 after creation")
	}

	if control.upperEndpoint != control.maxOutput {
		t.Errorf("expected upper endpoint to be %d but was %d", control.maxOutput, control.upperEndpoint)
	}
}

func Test8BitInput(t *testing.T) {
	control := NewPotControl(8, 12)

	testGetOutputParams(t, 255, 4095, UpperHalf, control)
	testGetOutputParams(t, 0, 4095/2, UpperHalf, control)

	testGetOutputParams(t, 255, 0, LowerHalf, control)
	testGetOutputParams(t, 0, 4095/2, LowerHalf, control)
}

func TestGetOutput(t *testing.T) {
	midpointInput := uint16(65535 / 2)
	midpointOutput := uint16(4095 / 2)

	control := NewPotControl(16, 12)

	checkUint16OutputValue(t, midpointOutput, control.GetOutputValue(midpointInput, FullWidth))

	checkUint16OutputValue(t, 4095, control.GetOutputValue(65535, FullWidth))

	checkUint16OutputValue(t, 0, control.GetOutputValue(0, FullWidth))

	checkUint16OutputValue(t, midpointOutput, control.GetOutputValue(0, UpperHalf))
	checkUint16OutputValue(t, midpointOutput, control.GetOutputValue(0, LowerHalf))

	checkUint16OutputValue(t, 4095, control.GetOutputValue(65535, UpperHalf))
	checkUint16OutputValue(t, 0, control.GetOutputValue(65535, LowerHalf))

	control.SetEndpoints(1000, 3900, 0)

	checkUint16OutputValue(t, midpointOutput, control.GetOutputValue(midpointInput, FullWidth))

	checkUint16OutputValue(t, 3900, control.GetOutputValue(65535, FullWidth))

	checkUint16OutputValue(t, 1000, control.GetOutputValue(0, FullWidth))

	checkUint16OutputValue(t, midpointOutput, control.GetOutputValue(0, UpperHalf))
	checkUint16OutputValue(t, midpointOutput, control.GetOutputValue(0, LowerHalf))

	checkUint16OutputValue(t, 3900, control.GetOutputValue(65535, UpperHalf))
	checkUint16OutputValue(t, 1000, control.GetOutputValue(65535, LowerHalf))

	control.SetEndpoints(1000, 3900, -12)

	checkUint16OutputValue(t, midpointOutput-12, control.GetOutputValue(midpointInput, FullWidth))
	checkUint16OutputValue(t, midpointOutput-12, control.GetOutputValue(0, LowerHalf))
	checkUint16OutputValue(t, midpointOutput-12, control.GetOutputValue(0, UpperHalf))

	control.SetEndpoints(12, 1800, -1000)

	checkUint16OutputValue(t, midpointOutput-1000, control.GetOutputValue(midpointInput, FullWidth))
	checkUint16OutputValue(t, midpointOutput-1000, control.GetOutputValue(0, LowerHalf))
	checkUint16OutputValue(t, midpointOutput-1000, control.GetOutputValue(0, UpperHalf))

	checkUint16OutputValue(t, 1423, control.GetOutputValue(49151, FullWidth))
	checkUint16OutputValue(t, 529, control.GetOutputValue(16383, FullWidth))

	control.Invert = true
	control.SetEndpoints(0, 4095, 0)

	invertedMidpoint := uint16(2048)
	checkUint16OutputValue(t, invertedMidpoint, control.GetOutputValue(0, UpperHalf))
	checkUint16OutputValue(t, invertedMidpoint, control.GetOutputValue(0, LowerHalf))
	checkUint16OutputValue(t, invertedMidpoint, control.GetOutputValue(midpointInput, FullWidth))

	checkUint16OutputValue(t, 0, control.GetOutputValue(65535, FullWidth))
	checkUint16OutputValue(t, 4095, control.GetOutputValue(0, FullWidth))

	checkUint16OutputValue(t, 0, control.GetOutputValue(65535, UpperHalf))
	checkUint16OutputValue(t, 4095, control.GetOutputValue(65535, LowerHalf))

	control.SetEndpoints(193, 2890, 0)
}

func checkUint16OutputValue(t *testing.T, expectedOutput, actualOutput uint16) {
	if expectedOutput != actualOutput {
		t.Errorf("expected output value to be %d but it was %d", expectedOutput, actualOutput)
	}
}

func testGetOutputParams(t *testing.T, inputValue, expectedOutput uint16, controlRange Range, c *PotControl) {
	actualOutput := c.GetOutputValue(inputValue, controlRange)

	if expectedOutput != actualOutput {
		t.Errorf("expected output value to be %d but it was %d. input: %d controlrange: %d control: %+ v", expectedOutput, actualOutput, inputValue, controlRange, c)
	}
}
