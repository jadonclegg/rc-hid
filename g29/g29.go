package g29

import (
	"encoding/binary"
	"math"
	"rc-hid/g29/button"

	"github.com/sstallion/go-hid"
)

const SteeringOffset = 4
const GasOffset = SteeringOffset + 2
const BrakeOffset = GasOffset + 1
const ClutchOffset = BrakeOffset + 1

const VendorID = 0x046d
const ProductID = 0xc24F
const Neutral = 32768

type WheelInput struct {
	Throttle uint16
	Brake    uint16
	Clutch   uint16
	Gear     byte
	Steering uint16
	raw      []byte
}

type buttonListener struct {
	lastVal bool
	action  ButtonAction
	btn     button.Button
}

type Wheel struct {
	Input           *WheelInput
	hidDevice       *hid.Device
	AutoCenterScale float64
	buttonListeners []*buttonListener
}

func NewWheel() (*Wheel, error) {
	wheel := &Wheel{}
	wheel.Input = &WheelInput{
		Brake:    0,
		Throttle: 0,
		Gear:     0,
		Steering: Neutral,
	}

	wheel.buttonListeners = make([]*buttonListener, 0)

	device, err := hid.OpenFirst(VendorID, ProductID)
	if err != nil {
		return nil, err
	}

	wheel.hidDevice = device
	return wheel, nil
}

type ButtonAction func(w *Wheel)

func (w *Wheel) OnButtonPress(btn button.Button, action ButtonAction) {
	listener := &buttonListener{
		lastVal: false,
		action:  action,
		btn:     btn,
	}

	w.buttonListeners = append(w.buttonListeners, listener)
}

func (wheel *Wheel) GetInput() {
	wheel.Input.raw = make([]byte, 12)
	wheel.Input.raw[0] = 0x01
	buff := wheel.Input.raw

	for {
		read, err := wheel.hidDevice.Read(buff)
		if err != nil {
			panic(err)
		}

		if read > 0 {
			wheel.Input.Steering = getUint16FromByteArray(buff, SteeringOffset, false)
			wheel.Input.Throttle = getUint8FromByteArrayAsUint16(buff, GasOffset, true)
			wheel.Input.Brake = getUint8FromByteArrayAsUint16(buff, BrakeOffset, true)
			wheel.Input.Clutch = getUint8FromByteArrayAsUint16(buff, ClutchOffset, true)
			wheel.Input.Gear = buff[button.FirstGear.Offset]

			for _, listener := range wheel.buttonListeners {
				if wheel.ButtonPressed(listener.btn) {
					if !listener.lastVal {
						listener.lastVal = true
						listener.action(wheel)
					}
				} else {
					listener.lastVal = false
				}
			}

			const maxForce = 60

			if wheel.AutoCenterScale > 0 {
				force := uint8(127)

				center := uint16(65535 / 2)
				distFromCenter := float64(wheel.Input.Steering) - float64(center)
				power := float64(0)

				if wheel.Input.Gear == button.ThirdGear.ByteVal {
					ratio := math.Min(1, math.Max(-1, float64(distFromCenter)/float64(center)*4))
					throttle := float64(wheel.Input.Throttle) * .5
					power = math.Max(-maxForce, math.Min(maxForce, float64(throttle)*ratio*wheel.AutoCenterScale))
				}

				force = 127 + uint8(power)

				wheel.SetConstantForce(force)
			}
		}
	}
}

func getUint16FromByteArray(buff []byte, offset int, invert bool) uint16 {
	byte_0 := buff[offset]
	byte_1 := buff[offset+1]

	value := uint16(0)
	value += (uint16(byte_1) << 8)
	value += uint16(byte_0)

	if invert {
		return 65535 - value
	}

	return value
}

func getUint8FromByteArrayAsUint16(buff []byte, offset int, invert bool) uint16 {
	value := uint16(buff[offset])
	if invert {
		return 255 - value
	}

	return value
}

func (w *Wheel) SetRange(degrees uint16) {
	if degrees < 40 {
		degrees = 40
	}

	if degrees > 900 {
		degrees = 900
	}

	degreeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(degreeBytes, degrees)

	bytes := []byte{0x00, 0xf8, 0x81, degreeBytes[0], degreeBytes[1], 0x00, 0x00, 0x00}
	w.hidDevice.Write(bytes)
}

func (w *Wheel) SetFrictionForce(value uint8) {
	bytes := []byte{0x00, 0x21, 0x02, byte(value), 0x00, byte(value), 0x00, 0x00}
	w.hidDevice.Write(bytes)
}

func (w *Wheel) SetConstantForce(value uint8) {
	bytes := []byte{0x00, 0x11, 0x00, byte(value), 0x00, 0x00, 0x00, 0x00}
	w.hidDevice.Write(bytes)
}

func (w *Wheel) ButtonPressed(button button.Button) bool {
	val := w.Input.raw[button.Offset]
	if val&button.ByteVal == button.ByteVal {
		return true
	}

	return false
}
