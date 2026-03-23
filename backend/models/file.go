package models

import (
	"encoding/json"
	"time"
)

type FileInfo struct {
	ID           string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	OriginalName string    `json:"original_name"`
	StorageName  string    `json:"storage_name"`
	StoragePath  string    `json:"storage_path"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	Extension    string    `json:"extension"`
	Hash         string    `json:"hash,omitempty"`
	URL          string    `json:"url"`
}

func (f *FileInfo) MarshalJSON() ([]byte, error) {
	type Alias FileInfo
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		CreatedAt: f.CreatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(f),
	})
}

func (f *FileInfo) UnmarshalJSON(data []byte) error {
	type Alias FileInfo
	aux := &struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.CreatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, aux.CreatedAt)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05", aux.CreatedAt)
			if err != nil {
				f.CreatedAt = time.Now()
			}
		}
		f.CreatedAt = parsedTime
	} else {
		f.CreatedAt = time.Now()
	}

	return nil
}
