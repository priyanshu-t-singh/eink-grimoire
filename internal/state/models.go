package state

import "time"

// PageType defines the current view tier in the Kavita hierarchy
type PageType string

const (
	PageLibrary  PageType = "library"
	PageSeries   PageType = "series"
	PageVolumes  PageType = "volumes"
	PageChapters PageType = "chapters"
	PageReader   PageType = "reader"
)

type StackFrame struct {
	Type PageType `json:"type"`

	// Kavita Context IDs (0 if not applicable to the current PageType)
	SeriesID  int `json:"series_id,omitempty"`
	VolumeID  int `json:"volume_id,omitempty"`
	ChapterID int `json:"chapter_id,omitempty"`

	// Cursor represents the currently highlighted index in a list (Library/Series/Volumes)
	// OR the current paginated page index in the Reader.
	Cursor int `json:"cursor"`
}

// DeviceState represents the full persisted state for a single e-ink device.
type DeviceState struct {
	DeviceID      string       `json:"device_id"`
	NavStack      []StackFrame `json:"nav_stack"`
	LastFrameHash string       `json:"last_frame_hash"` // Used to skip rendering if state hasn't changed
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (s *DeviceState) Top() *StackFrame {
	if len(s.NavStack) == 0 {
		return &StackFrame{Type: PageLibrary, Cursor: 0}
	}
	return &s.NavStack[len(s.NavStack)-1]
}

func (s *DeviceState) Pop() {
	if len(s.NavStack) > 1 {
		s.NavStack = s.NavStack[:len(s.NavStack)-1]
	} else {
		// Empty stack defaults to Library root
		s.NavStack = []StackFrame{{Type: PageLibrary, Cursor: 0}}
	}
}

func (s *DeviceState) Push(frame StackFrame) {
	s.NavStack = append(s.NavStack, frame)
}
