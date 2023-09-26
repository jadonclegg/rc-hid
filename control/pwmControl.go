package control

// ~half of max of uint16
const Neutral = 32768

type ControlRange int

const (
	FULL_WIDTH ControlRange = 0
	LOWER_HALF ControlRange = 1
	UPPER_HALF ControlRange = 2
)

type PwmControl struct {
	Trim   int
	Invert bool
	endpointInfo
}

func NewPwmControl(resolution uint16) *PwmControl {
	newControl := &PwmControl{
		Trim: 0,
	}

	newControl.endpointInfo = *newEndpointInfo(resolution)

	return newControl
}

func (c *PwmControl) getEnpointInfo() endpointInfo {
	return c.endpointInfo
}

func (c *PwmControl) GetDutyCycle(input uint16, controlRange ControlRange) uint16 {
	adjustedValue := adjustValueForEndpoints(input, c.endpointInfo)
	adjustedValue = c.trim(adjustedValue)

	// fmt.Println(adjustedValue)
	if c.Invert {
		adjustedValue = c.maxOutput - adjustedValue
	}

	if controlRange == FULL_WIDTH {
		return uint16(int64(DUTY_CYCLE_RANGE)*int64(adjustedValue)/int64(c.maxOutput)) + MIN_DUTY_CYCLE
	}

	dutyVal := uint16(int64(HALF_DUTY_CYCLE_RANGE) * int64(adjustedValue) / int64(c.maxOutput))

	if controlRange == LOWER_HALF {
		return (HALF_DUTY_CYCLE_RANGE - dutyVal) + MIN_DUTY_CYCLE
	}

	return dutyVal + HALF_DUTY_CYCLE_RANGE + MIN_DUTY_CYCLE
}

func (c *PwmControl) trim(value uint16) uint16 {
	return trimValue(value, c)
}

func (c *PwmControl) getTrim() int {
	return c.Trim
}
