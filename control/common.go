package control

const FREQUENCY = 50 // hz
const MIN_US = 500   // .5ms pulse
const MAX_US = 2500  // 2.5ms pulse

// 20 ms
const PERIOD = 1000000 / FREQUENCY // microseconds in each period

const MAX_DUTY_CYCLE = 65535 * MAX_US / PERIOD
const MIN_DUTY_CYCLE = 65535 * MIN_US / PERIOD
const DUTY_CYCLE_RANGE = MAX_DUTY_CYCLE - MIN_DUTY_CYCLE
const HALF_DUTY_CYCLE_RANGE = DUTY_CYCLE_RANGE / 2

type trimmableControl interface {
	getTrim() int
}

func trimValue(value uint16, c trimmableControl) uint16 {
	trim := c.getTrim()

	if trim > 0 {
		return value + uint16(trim)
	}

	if trim < 0 {
		return value - uint16(trim*-1)
	}

	return value
}

func min(a, b uint16) uint16 {
	if a < b {
		return a
	}

	return b
}

func max(a, b uint16) uint16 {
	if a > b {
		return a
	}

	return b
}
