package app

func (a *App) GetNewestRecord() *Record {
	a.Lock()
	defer a.Unlock()

	if len(a.records) == 0 {
		return nil
	}
	rec := a.records[0]
	for i := 1; i < len(a.records); i++ {
		if rec.Start.Before(a.records[i].Start) {
			rec = a.records[i]
		}
	}
	return &rec
}

func (a *App) AddRecord(rec Record) error {
	a.Lock()
	defer a.Unlock()

	// TODO Validate fields and conflicts, return error.

	a.lastID++
	rec.ID = a.lastID
	a.records = append(a.records, rec)
	a.emitRecordAdded(EventRecordAdded{rec})
	return nil
}
