package button

// Button is used to map which buttons are where on the g29 wheel, and what part of the input byte array they are.
type Button struct {
	Flag   byte
	Offset int
}

// Buttons
var (
	Cross       Button = Button{0x10, 0}
	Square      Button = Button{0x20, 0}
	Circle      Button = Button{0x40, 0}
	Triangle    Button = Button{0x80, 0}
	Upshift     Button = Button{0x01, 1}
	Downshift   Button = Button{0x02, 1}
	R2          Button = Button{0x04, 1}
	L2          Button = Button{0x08, 1}
	Share       Button = Button{0x10, 1}
	Option      Button = Button{0x20, 1}
	R3          Button = Button{0x40, 1}
	L3          Button = Button{0x80, 1}
	Plus        Button = Button{0x80, 2}
	Minus       Button = Button{0x01, 3}
	SpinRight   Button = Button{0x02, 3}
	SpinLeft    Button = Button{0x04, 3}
	SpinButton  Button = Button{0x08, 3}
	Playstation Button = Button{0xF0, 3}
)

// Shifter gears
var (
	FirstGear  Button = Button{0x01, 2}
	SecondGear Button = Button{0x02, 2}
	ThirdGear  Button = Button{0x04, 2}
	FourthGear Button = Button{0x08, 2}
	FifthGear  Button = Button{0x10, 2}
	SixthGear  Button = Button{0x20, 2}
	Reverse    Button = Button{0x40, 2}
)

// D-Pad
var (
	Top         Button = Button{0x00, 0}
	TopRight    Button = Button{0x01, 0}
	Right       Button = Button{0x02, 0}
	BottomRight Button = Button{0x03, 0}
	Bottom      Button = Button{0x04, 0}
	BottomLeft  Button = Button{0x05, 0}
	Left        Button = Button{0x06, 0}
	TopLeft     Button = Button{0x07, 0}
	DPadNeutral Button = Button{0x08, 0}
)
