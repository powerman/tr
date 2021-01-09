package comp

import (
	"errors"
	"fmt"
	"time"

	"github.com/powerman/tr/web/app/internal/app"
)

const (
	day         = 24 * time.Hour
	step        = 15 * time.Minute
	timeFmt     = "15:04"
	datetimeFmt = "2006-01-02 " + timeFmt
)

type Record struct {
	Date     string   `vugu:"data"` // "yyyy-mm-dd" or "".
	TimeFrom string   `vugu:"data"` // "hh:mm" or "".
	TimeTo   string   `vugu:"data"` // "hh:mm" or "".
	Activity string   `vugu:"data"`
	Actors   []string `vugu:"data"`
	Customer string   `vugu:"data"`
	Details  string   `vugu:"data"`
	Duration string   // "h:mm", computed.
}

func ToRecord(v app.Record) Record {
	return Record{
		Date:     v.Start.Format("2006-01-02"),
		TimeFrom: v.Start.Format("15:04"),
		TimeTo:   v.Start.Add(v.Duration).Format("15:04"),
		Activity: v.Activity,
		Actors:   v.Actors,
		Customer: v.Customer,
		Details:  v.Details,
		Duration: fmt.Sprintf("%d:%02d", int(v.Duration.Hours()), int(v.Duration.Minutes())%60),
	}
}

func (v Record) ToApp() (*app.Record, error) {
	m := &app.Record{
		Activity: v.Activity,
		Actors:   v.Actors,
		Customer: v.Customer,
		Details:  v.Details,
	}
	end, err := time.Parse(datetimeFmt, v.Date+" "+v.TimeTo)
	if err == nil {
		m.Start, err = time.Parse(datetimeFmt, v.Date+" "+v.TimeFrom)
	}
	if err == nil {
		m.Duration = end.Add(day).Sub(m.Start).Truncate(step) % day
		if m.Duration <= 0 {
			err = errors.New("no duration")
		}
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}
