package control

// Range is used for controlling the range of the output that you want to use from your input.
type Range int

// Control Ranges for outputs. Full width maps the input to the full range of the output, lowerhalf maps it to the lower half of the range, starting at the middle, and upper half does the upper half of the range, also starting at the middle.
const (
	FullWidth Range = 0
	LowerHalf Range = 1
	UpperHalf Range = 2
)

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
