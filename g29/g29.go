package g29

import (
	"encoding/binary"
	"math"
	"rc-hid/g29/button"

	"github.com/sstallion/go-hid"
)

const steeringOffset = 4
const gasOffset = steeringOffset + 2
const brakeOffset = gasOffset + 1
const clutchOffset = brakeOffset + 1

const vendorID = 0x046d
const productID = 0xc24F
const neutral = 32768

// WheelInput contains all of the input information coming from the wheel.
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

// Wheel contains information about the g29 wheel, and contains the WheelInput object for getting input values.
// AutoCenterScale is used to control the strength of the auto centering feature.
type Wheel struct {
	Input           *WheelInput
	hidDevice       *hid.Device
	AutoCenterScale float64
	buttonListeners []*buttonListener
}

// NewWheel is used to open and initialize the steering wheel. Returns an error if something went wrong.
func NewWheel() (*Wheel, error) {
	wheel := &Wheel{}
	wheel.Input = &WheelInput{
		Brake:    0,
		Throttle: 0,
		Gear:     0,
		Steering: neutral,
	}

	wheel.buttonListeners = make([]*buttonListener, 0)

	device, err := hid.OpenFirst(vendorID, productID)
	if err != nil {
		return nil, err
	}

	wheel.hidDevice = device
	return wheel, nil
}

// ButtonAction is a function that is called on a button press, and just passes the wheel object into the function that had the event.
type ButtonAction func(w *Wheel)

// OnButtonPress calls the passed in ButtonAction whenever the specified button is pressed.
func (w *Wheel) OnButtonPress(btn button.Button, action ButtonAction) {
	listener := &buttonListener{
		lastVal: false,
		action:  action,
		btn:     btn,
	}

	w.buttonListeners = append(w.buttonListeners, listener)
}

// GetInput runs forever and constantly reads input from the device, and stores the information in the wheel.Input field.
func (w *Wheel) GetInput() {
	w.Input.raw = make([]byte, 12)
	w.Input.raw[0] = 0x01
	buff := w.Input.raw

	for {
		read, err := w.hidDevice.Read(buff)
		if err != nil {
			panic(err)
		}

		if read > 0 {
			w.Input.Steering = getUint16FromByteArray(buff, steeringOffset, false)
			w.Input.Throttle = getUint8FromByteArrayAsUint16(buff, gasOffset, true)
			w.Input.Brake = getUint8FromByteArrayAsUint16(buff, brakeOffset, true)
			w.Input.Clutch = getUint8FromByteArrayAsUint16(buff, clutchOffset, true)
			w.Input.Gear = buff[button.FirstGear.Offset]

			for _, listener := range w.buttonListeners {
				if w.ButtonPressed(listener.btn) {
					if !listener.lastVal {
						listener.lastVal = true
						listener.action(w)
					}
				} else {
					listener.lastVal = false
				}
			}

			const maxForce = 60

			if w.AutoCenterScale > 0 {
				force := uint8(127)

				center := uint16(65535 / 2)
				distFromCenter := float64(w.Input.Steering) - float64(center)
				power := float64(0)

				if w.Input.Gear == button.ThirdGear.Flag {
					ratio := math.Min(1, math.Max(-1, float64(distFromCenter)/float64(center)*4))
					throttle := float64(w.Input.Throttle) * .5
					power = math.Max(-maxForce, math.Min(maxForce, float64(throttle)*ratio*w.AutoCenterScale))
				}

				force = 127 + uint8(power)

				w.SetConstantForce(force)
			}
		}
	}
}

func getUint16FromByteArray(buff []byte, offset int, invert bool) uint16 {
	byte0 := buff[offset]
	byte1 := buff[offset+1]

	value := uint16(0)
	value += (uint16(byte1) << 8)
	value += uint16(byte0)

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

// SetRange lets you specify how many degrees of rotation you want the wheel to have, anywhere from 40 to 900 degrees.
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

// SetFrictionForce enables the friction force feedback of the wheel, to enable a constant resistance to turning the wheel. Values are 0-F
func (w *Wheel) SetFrictionForce(value uint8) {
	bytes := []byte{0x00, 0x21, 0x02, byte(value), 0x00, byte(value), 0x00, 0x00}
	w.hidDevice.Write(bytes)
}

// SetConstantForce tells the wheel to enable a constant force feedback. Values between 0-127 turn the wheel left, with 0 being the strongest feedback,
// while values from 128-255 turn the wheel to the right, with 255 being the strongest feedback.
func (w *Wheel) SetConstantForce(value uint8) {
	bytes := []byte{0x00, 0x11, 0x00, byte(value), 0x00, 0x00, 0x00, 0x00}
	w.hidDevice.Write(bytes)
}

// ButtonPressed lets you check if a specific button is being pressed currently.
func (w *Wheel) ButtonPressed(button button.Button) bool {
	val := w.Input.raw[button.Offset]
	if val&button.Flag == button.Flag {
		return true
	}

	return false
}
