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
			// Lazy initialization — matches state.NewDeviceState's root page exactly.
			return state.NewDeviceState(deviceID), nil
		}
		return nil, err
	}

	var stack []state.Page
	if err := json.Unmarshal(navStackJSON, &stack); err != nil {
		return nil, err
	}

	ds := &state.DeviceState{
		DeviceID:  deviceID,
		Stack:     stack,
		UpdatedAt: updatedAt,
	}
	if lastFrameHash.Valid {
		ds.LastFrameHash = &lastFrameHash.String
	}
	return ds, nil
}

// SaveDeviceState upserts the current navigation stack and frame hash.
func (r *Repository) SaveDeviceState(s *state.DeviceState) error {
	navStackJSON, err := json.Marshal(s.Stack)
	if err != nil {
		return err
	}

	var lastFrameHash sql.NullString
	if s.LastFrameHash != nil {
		lastFrameHash = sql.NullString{String: *s.LastFrameHash, Valid: true}
	}

	query := `
		INSERT INTO device_states (device_id, nav_stack, last_frame_hash, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(device_id) DO UPDATE SET
			nav_stack = excluded.nav_stack,
			last_frame_hash = excluded.last_frame_hash,
			updated_at = CURRENT_TIMESTAMP;
	`

	_, err = r.db.Exec(query, s.DeviceID, navStackJSON, lastFrameHash)
	return err
}
