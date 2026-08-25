package state

import "time"

// PageType defines the current view tier in the Kavita hierarchy
type PageType string

const (
	PageLibrary  PageType = "Library"
	PageSeries   PageType = "Series"
	PageBookList PageType = "BookList"
	PageReader   PageType = "Reader"
)

// Params: immutable identifiers needed to refetch content.
// State: mutable cursor/scroll/sub-page position, changes in place.
type Page struct {
	Type   PageType          `json:"type"`
	Params map[string]string `json:"params,omitempty"`
	State  map[string]int    `json:"state"`
}

// DeviceState represents the full persisted state for a single e-ink device.
type DeviceState struct {
	DeviceID      string    `json:"device_id"`
	Stack         []Page    `json:"stack"`
	LastFrameHash *string   `json:"last_frame_hash,omitempty"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewDeviceState(deviceID string) *DeviceState {
	return &DeviceState{
		DeviceID: deviceID,
		Stack: []Page{
			{Type: PageLibrary, State: map[string]int{"cursor": 0, "scroll": 0}},
		},
		UpdatedAt: time.Now(),
	}
}

func (s *DeviceState) Top() *Page {
	if len(s.Stack) == 0 {
		s.Stack = append(s.Stack, Page{Type: PageLibrary, State: map[string]int{"cursor": 0, "scroll": 0}})
	}
	return &s.Stack[len(s.Stack)-1]
}

func (s *DeviceState) Pop() {
	if len(s.Stack) > 1 {
		s.Stack = s.Stack[:len(s.Stack)-1]
	} else {
		// Empty stack defaults to Library root
		s.Stack = []Page{{Type: PageLibrary, State: map[string]int{"cursor": 0, "scroll": 0}}}
	}
}

func (s *DeviceState) Push(page Page) {
	s.Stack = append(s.Stack, page)
}
