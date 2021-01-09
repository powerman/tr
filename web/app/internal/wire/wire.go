package wire

import "github.com/powerman/tr/web/app/internal/app"

type (
	ApplWirer interface{ WireAppl(app.Appl) }
	Appl      struct{ app.Appl }
)

func (vr *Appl) WireAppl(v app.Appl) { vr.Appl = v }
