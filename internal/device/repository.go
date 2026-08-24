package device

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"le-grimoire/internal/state"
)

type Repository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetDeviceAuthHash retrieves the stored API key hash for a given device.
func (r *Repository) GetDeviceAuthHash(deviceID string) (string, error) {
	var hash string
	err := r.db.QueryRow("SELECT api_key_hash FROM devices WHERE device_id = ?", deviceID).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("device not registered")
		}
		return "", err
	}
	return hash, nil
}

// GetDeviceState retrieves the current navigation back-stack for a device.
func (r *Repository) GetDeviceState(deviceID string) (*state.DeviceState, error) {
	var navStackJSON []byte
	var lastFrameHash sql.NullString
	var updatedAt time.Time

	query := `SELECT nav_stack, last_frame_hash, updated_at FROM device_states WHERE device_id = ?`
	err := r.db.QueryRow(query, deviceID).Scan(&navStackJSON, &lastFrameHash, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Lazy Initialization
			return &state.DeviceState{
				DeviceID: deviceID,
				NavStack: []state.StackFrame{
					{Type: state.PageLibrary, Cursor: 0},
				},
				LastFrameHash: "",
				UpdatedAt:     time.Now(),
			}, nil
		}
		return nil, err
	}

	var navStack []state.StackFrame
	if err := json.Unmarshal(navStackJSON, &navStack); err != nil {
		return nil, err
	}

	return &state.DeviceState{
		DeviceID:      deviceID,
		NavStack:      navStack,
		LastFrameHash: lastFrameHash.String,
		UpdatedAt:     updatedAt,
	}, nil
}

// SaveDeviceState upserts the current navigation stack and frame hash.
func (r *Repository) SaveDeviceState(s *state.DeviceState) error {
	navStackJSON, err := json.Marshal(s.NavStack)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO device_states (device_id, nav_stack, last_frame_hash, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(device_id) DO UPDATE SET
			nav_stack = excluded.nav_stack,
			last_frame_hash = excluded.last_frame_hash,
			updated_at = CURRENT_TIMESTAMP;
	`

	_, err = r.db.Exec(query, s.DeviceID, navStackJSON, s.LastFrameHash)
	return err
}
