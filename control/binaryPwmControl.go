package control

type BinaryPwmControl struct {
	lowerEndpoint uint16
	upperEndpoint uint16
	Inverted      bool
}

func NewBinaryPwmControl() *BinaryPwmControl {
	control := &BinaryPwmControl{
		lowerEndpoint: 0,
		upperEndpoint: 65535,
		Inverted:      false,
	}

	return control
}

func (c *BinaryPwmControl) GetDutyCycle(on bool) uint16 {
	value := uint16(0)

	if on {
		value = 65535
	} else {
		value = 0
	}

	if c.Inverted {
		value = 65535 - value
	}

	adjustedValue := min(c.upperEndpoint, max(value, c.lowerEndpoint))

	return uint16(int64(DUTY_CYCLE_RANGE)*int64(adjustedValue)/int64(65535)) + MIN_DUTY_CYCLE
}
