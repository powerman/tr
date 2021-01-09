package ui

import (
	"sort"

	"github.com/vugu/vugu"

	"github.com/powerman/tr/web/app/internal/app"
	"github.com/powerman/tr/web/app/internal/ui/comp"
	"github.com/powerman/tr/web/app/internal/wire"
)

type PageReport struct {
	wire.Appl
	Records []*comp.Record `vugu:"data"` // Loaded TR records in reverse order (last one is at [0]).
}

func (c *PageReport) Init(ctx vugu.InitCtx) {
	c.Appl.SubscribeRecordAdded(c, ctx)
}

func (c *PageReport) Destroy(ctx vugu.DestroyCtx) {
	c.Appl.UnsubscribeRecordAdded(c)
}

func (c *PageReport) OnRecordAdded(ev app.EventRecordAdded) bool {
	r := comp.ToRecord(ev.Data)
	pos := sort.Search(len(c.Records), func(i int) bool {
		return r.Date+r.TimeFrom > c.Records[i].Date+c.Records[i].TimeFrom
	})
	c.Records = append(c.Records, nil)
	copy(c.Records[pos+1:], c.Records[pos:])
	c.Records[pos] = &r
	return true
}
