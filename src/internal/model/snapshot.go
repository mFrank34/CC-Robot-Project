package model

import (
	"encoding/json"

	"Robot-Project/internal/util"
)

type snapshot struct {
	Bots   map[string]Bot     `json:"bots"`
	Latest map[string]Message `json:"latest"`
	Status map[string]Status  `json:"status"`
}

func (s *Store) Save(path string) error {
	s.mu.RLock()
	snap := snapshot{Bots: s.bots, Latest: s.latest, Status: s.status}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return util.WriteJSONAtomic(path, data)
}

func Load(path string) (*Store, error) {
	data, found, err := util.ReadFileOrEmpty(path)
	if err != nil {
		return nil, err
	}
	if !found {
		return NewStore(), nil
	}

	var snap snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	if snap.Bots == nil {
		snap.Bots = make(map[string]Bot)
	}
	if snap.Latest == nil {
		snap.Latest = make(map[string]Message)
	}
	if snap.Status == nil {
		snap.Status = make(map[string]Status)
	}

	return &Store{bots: snap.Bots, latest: snap.Latest, status: snap.Status}, nil
}
