package main

import (
	"encoding/json"
)

// EventHandler _
type EventHandler func(*Event)

// Event _
type Event struct {
	Name string      `json:"event"`
	Data interface{} `json:"data"`
}

// NewEvent _
func NewEvent(rawData []byte) (*Event, error) {
	event := new(Event)
	err := json.Unmarshal(rawData, event)
	return event, err
}

// Raw _
func (e *Event) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}
