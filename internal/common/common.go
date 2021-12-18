package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	// NameGauge is the string name of the gauge type
	NameGauge = "gauge"
	// NameCounter is the string name of the counter type
	NameCounter = "counter"
)

// Metrics is the struct to use for metrics updates
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

// ComputeHash calculates hash for metrics
func (m Metrics) ComputeHash(key string) (*[sha256.Size]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("no key")
	}
	if m.ID == "" {
		return nil, fmt.Errorf("empty ID field")
	}
	toHash := ""
	if m.MType == NameGauge {
		if m.Value == nil {
			return nil, fmt.Errorf("no value")
		}
		toHash = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value) + key
	}
	if m.MType == NameCounter {
		if m.Delta == nil {
			return nil, fmt.Errorf("no delta")
		}
		toHash = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta) + key
	}

	hash := sha256.Sum256([]byte(toHash))
	return &hash, nil
}

// StoreHash stores hash to Metrics struct
func (m *Metrics) StoreHash(key string) error {
	if key == "" {
		return nil
	}
	h, err := m.ComputeHash(key)
	if err != nil {
		return nil
	}
	m.Hash = hex.EncodeToString(h[:])
	return nil
}

// CheckHash checks hash stored into Metrics struct
func (m Metrics) CheckHash(key string) error {
	if key == "" {
		return nil
	}
	h, err := m.ComputeHash(key)
	if err != nil {
		return nil
	}
	hashStr := hex.EncodeToString(h[:])
	if m.Hash != hashStr {
		return fmt.Errorf("hash value incorrect")
	}
	return nil
}
