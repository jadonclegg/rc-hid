package control

import "testing"

func TestCreateEndpointInfo(t *testing.T) {
	info := newEndpointInfo(16)
	if info.maxOutput != 65535 {
		t.Errorf("expected max output to be %d but it was %d", 65535, info.maxOutput)
	}

	if info.lowerEndpoint != 0 {
		t.Errorf("expected lower endpoint to initialize to 0 but it was %d", info.lowerEndpoint)
	}

	if info.upperEndpoint != info.maxOutput {
		t.Errorf("expected upper endpoint to initialize to %d but it was %d", info.maxOutput, info.upperEndpoint)
	}
}

func TestSetEndpoints(t *testing.T) {
	info := newEndpointInfo(16)
	info.SetEndpoints(12000, 45000, 0)

	if info.upperEndpoint != 45000 {
		t.Errorf("failed to set upper endpoint")
	}

	if info.lowerEndpoint != 12000 {
		t.Errorf("failed to set lower endpoint")
	}
}
