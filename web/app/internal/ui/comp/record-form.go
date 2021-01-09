package comp

import (
	"reflect"
	"time"

	"github.com/vugu/vugu"

	"github.com/powerman/tr/web/app/internal/app"
	"github.com/powerman/tr/web/app/internal/wire"
)

type RecordForm struct {
	wire.Appl
	ID         int    // 0 means "Add" form, other means "Edit" form.
	Record     Record `vugu:"data"`
	origRecord Record
}

func (c *RecordForm) Init(ctx vugu.InitCtx) {
	c.Appl.SubscribeRecordAdded(c, ctx)
}

func (c *RecordForm) Destroy(ctx vugu.DestroyCtx) {
	c.Appl.UnsubscribeRecordAdded(c)
}

func (c *RecordForm) Compute(ctx vugu.ComputeCtx) {
	if c.Record.Date == "" {
		c.reset()
	}
}

func (c *RecordForm) OnRecordAdded(_ app.EventRecordAdded) bool { return c.updateOrigRecord() }

func (c *RecordForm) updateOrigRecord() bool {
	if c.ID != 0 {
		return false
	}
	newest := c.Appl.GetNewestRecord()
	if newest == nil {
		return false
	}
	rec := ToRecord(*newest)
	if reflect.DeepEqual(c.origRecord, rec) {
		return false
	}
	reset := c.Record.Details == "" &&
		c.origRecord.Activity == c.Record.Activity &&
		reflect.DeepEqual(c.origRecord.Actors, c.Record.Actors) &&
		c.origRecord.Customer == c.Record.Customer
	c.origRecord = rec
	if reset {
		c.reset()
	}
	return reset
}

// In Firefox <input type=time step=N> does not enforce step, so let's do
// this manually.
func (c *RecordForm) fixStepFirefox(field *string) {
	t, err := time.Parse(timeFmt, *field)
	if err == nil {
		*field = t.Round(step).Format(timeFmt)
	}
}

func (c *RecordForm) add() {
	r, err := c.Record.ToApp()
	if err == nil {
		err = c.Appl.AddRecord(*r)
	}
	if err == nil {
		c.reset()
	}
}

func (c *RecordForm) reset() {
	c.Record = c.origRecord
	if c.ID == 0 {
		now := time.Now().Round(step)
		c.Record.Date = now.Format("2006-01-02")
		c.Record.TimeFrom = now.Format("15:04")
		c.Record.TimeTo = now.Add(step).Format("15:04")
		c.Record.Details = ""
	}
}
