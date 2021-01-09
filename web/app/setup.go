package main

import (
	"github.com/vugu/vugu"

	"github.com/powerman/tr/web/app/internal/app"
	"github.com/powerman/tr/web/app/internal/wire"
)

func vuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv) vugu.Builder {
	appl := app.New()

	buildEnv.SetWireFunc(func(b vugu.Builder) {
		if c, ok := b.(wire.ApplWirer); ok {
			c.WireAppl(appl)
		}
	})

	ret := &Root{}
	buildEnv.WireComponent(ret)

	appl.Initialize()

	return ret
}
