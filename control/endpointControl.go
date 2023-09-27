package control

import "math"

type endpointInfo struct {
	lowerEndpoint  uint16
	upperEndpoint  uint16
	ratio          float64
	upperRatio     float64
	lowerRatio     float64
	offset         uint16
	halfMaxOutput  float64
	maxOutput      uint16
	trim           int
	outputMidPoint float64
}

func newEndpointInfo(resolution uint16) *endpointInfo {
	max := math.Pow(2.0, float64(resolution)) - 1
	halfMax := max / float64(2)

	info := &endpointInfo{
		maxOutput:     uint16(max),
		halfMaxOutput: halfMax,
	}

	info.SetEndpoints(0, info.maxOutput, 0)

	return info
}

// SetEndpoints allows you to specify the min, and max output values that should be used, and lets you adjust where the middle of those values is with trim.
func (info *endpointInfo) SetEndpoints(lower, upper uint16, trim int) {
	info.lowerEndpoint = lower
	info.upperEndpoint = upper

	info.outputMidPoint = info.halfMaxOutput + float64(trim)

	usableRange := info.upperEndpoint - info.lowerEndpoint
	info.ratio = float64(usableRange) / float64(info.maxOutput)
	info.offset = (info.maxOutput - usableRange) / 2

	info.upperRatio = (float64(info.upperEndpoint) - info.outputMidPoint) / info.halfMaxOutput
	info.lowerRatio = (info.outputMidPoint - float64(info.lowerEndpoint)) / info.halfMaxOutput
}

func adjustValueForEndpoints(value uint16, info endpointInfo) uint16 {
	lowerEndpoint := info.lowerEndpoint
	upperEndpoint := info.upperEndpoint
	max := info.maxOutput

	if lowerEndpoint == 0 {
		return uint16(float64(value) * info.ratio)
	}

	if upperEndpoint == max {
		return uint16(float64(value)*info.ratio + float64(lowerEndpoint))
	}

	if value > uint16(info.halfMaxOutput) {
		return uint16((float64(value)-info.halfMaxOutput)*info.upperRatio + info.outputMidPoint)
	}

	return uint16(info.outputMidPoint - ((info.halfMaxOutput - float64(value)) * info.lowerRatio))
}
