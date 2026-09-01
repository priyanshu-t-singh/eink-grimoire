package port

import "time"

type Book struct {
	ID, Title, Author string
	TotalPages        int
}

type PageFormat int

const (
	FormatText PageFormat = iota
	FormatImage
)

type Page struct {
	BookID  string
	Number  int
	Content []byte
	Format  PageFormat
}

type Progress struct {
	BookID     string
	PageNumber int
	UpdatedAt  time.Time
}

type ColorMode int

const (
	ColorMono ColorMode = iota
	ColorANSI256
	ColorRGB
)

type Capabilities struct {
	Width, Height int
	Color         ColorMode
	Rotation      int
}

type ButtonEvent int

const (
	ButtonNext ButtonEvent = iota
	ButtonPrev
	ButtonMenu
)

// Event ties a raw button press to the session (device/tab/user) it came from.
type Event struct {
	SessionID string
	Button    ButtonEvent
}
