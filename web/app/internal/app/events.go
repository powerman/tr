package app

import "github.com/vugu/vugu"

type EventEnver interface { // Probably should be defined in vugu package.
	EventEnv() vugu.EventEnv
}

type (
	EventInitialized struct{}
	EventRecordAdded struct {
		Data Record
	}
)
