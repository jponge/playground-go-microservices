package model

import (
	"errors"
	"fmt"
	"sync"
)

type Database struct {
	entries map[string]float64
	mutex   sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{entries: map[string]float64{}, mutex: sync.RWMutex{}}
}

func (db *Database) Get(id string) (float64, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	if temp, found := db.entries[id]; found {
		return temp, found
	}
	return 0, false
}

func (db *Database) Put(id string, temp float64) {
	db.mutex.Lock()
	db.entries[id] = temp
	db.mutex.Unlock()
}

func (db *Database) AllTemperatureUpdates() []TemperatureUpdate {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	var entries []TemperatureUpdate
	for k, v := range db.entries {
		entries = append(entries, TemperatureUpdate{SensorId: k, Value: v})
	}
	return entries
}

func (db *Database) OneTemperatureUpdate(id string) (*TemperatureUpdate, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	if value, found := db.entries[id]; found {
		return &TemperatureUpdate{
			SensorId: id,
			Value:    value,
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("No entry for key %s", id))
}
