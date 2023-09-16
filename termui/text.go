package termui

type FontColor int

const (
	Green FontColor = iota
	Red
	Yellow
	White
	Blue
)

type FontJustification int

const (
	Left FontJustification = iota
	Right
	Center
)

type TextOptions func(*CustomText) error
type CustomText struct {
	text  string
	color FontColor
	x     int
	y     int
	width int
}

func Color(c FontColor) func(*CustomText) error {
	return func(t *CustomText) error {
		t.color = c
		return nil
	}
}

func Justify(j FontJustification) func(*CustomText) error {
	return func(t *CustomText) error {
		switch j {
		case Right:
			t.x -= len(t.text)
		case Center:
			t.x -= (len(t.text) / 2)
		case Left:
			// nothing to do here
		}

		return nil
	}
}

// Width is used to set the minimum width of the text. This is useful for
// clearing out old text that was longer than the new text
func Width(w int) func(*CustomText) error {
	return func(t *CustomText) error {
		t.width = w
		return nil
	}
}
