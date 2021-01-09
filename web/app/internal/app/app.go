package app

import (
	"sync"
	"time"
)

type Appl interface {
	// GetNewestRecord returns newest or nil if there are no records.
	GetNewestRecord() *Record
	AddRecord(Record) error
	TopicInitialized
	TopicRecordAdded
}

type Record struct {
	ID       int
	Start    time.Time
	Duration time.Duration // 15m…23h45m, step 15m.
	Activity string
	Actors   []string
	Customer string
	Details  string
}

type App struct {
	sync.Mutex
	lastID  int
	records []Record
	topicInitialized
	topicRecordAdded
}

func New() *App {
	return &App{}
}

func (a *App) Initialize() {
	a.loadFixture()
	a.emitInitialized(EventInitialized{})
}
