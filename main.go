package main

import (
	"bufio"
	"fmt"
	"rc-hid/control"
	"rc-hid/g29"
	g29button "rc-hid/g29/button"
	"time"

	"go.bug.st/serial"

	"github.com/sstallion/go-hid"
)

const SWA_ON = 0x80
const SWD_ON = 0x40

const SWB_HIGH = 0x20
const SWB_MID = 0x10

const SWC_HIGH = 0x08
const SWC_MID = 0x04

func main() {
	fmt.Println("RC HID")

	strControl := control.NewPotControl(16, 12)
	strControl.SetEndpoints(15, 3760, -217)
	strControl.Invert = true

	thrControl := control.NewPotControl(8, 12)
	thrControl.SetEndpoints(335, 2790, -122)
	thrControl.Invert = true

	vrControl := control.NewPotControl(16, 12)
	vrControl.SetEndpoints(659, 3455, 0)

	lightsOn := false

	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		panic(err)
	}

	go readSerialOutput(port)

	wheel, err := g29.NewWheel()
	if err != nil {
		panic(err)
	}

	wheel.SetRange(360)
	wheel.SetFrictionForce(2)
	wheel.AutoCenterScale = .5

	wheel.OnButtonPress(g29button.Triangle, func(w *g29.Wheel) {
		lightsOn = !lightsOn
	})

	go wheel.GetInput()

	for {
		time.Sleep(time.Millisecond * 15)

		var thrOutput uint16
		strOutput := strControl.GetOutputValue(wheel.Input.Steering, control.FullWidth)
		vrOutput := vrControl.GetOutputValue(wheel.Input.Steering, control.FullWidth)

		if wheel.Input.Gear == g29button.ThirdGear.Flag {
			thrOutput = thrControl.GetOutputValue(wheel.Input.Throttle, control.UpperHalf)
		} else if wheel.Input.Gear == g29button.FourthGear.Flag {
			thrOutput = thrControl.GetOutputValue(wheel.Input.Throttle, control.LowerHalf)
		} else {
			thrOutput = thrControl.GetOutputValue(0, control.UpperHalf)
		}

		var buttons uint16

		if lightsOn {
			buttons = buttons | SWB_HIGH
		}

		// Third spot is a placeholder for vr knob on transmitter once it's wired. I'm lazy.
		outputStr := fmt.Sprintf("%04X%04X%04X%02X\n", strOutput, thrOutput, vrOutput, buttons)
		port.Write([]byte(outputStr))
	}
}

func listHid() {
	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		fmt.Printf("%s: ID %04x:%04x %s %s\n",
			info.Path,
			info.VendorID,
			info.ProductID,
			info.MfrStr,
			info.ProductStr)

		return nil
	})
}

func readSerialOutput(port serial.Port) {
	scanner := bufio.NewScanner(port)

	for {
		if scanner.Scan() {

			fmt.Println("Got input: ", string(scanner.Bytes()))
		} else if scanner.Err() != nil {
			panic("Disconnected")
		}
	}
}
