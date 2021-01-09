package ui

// DO NOT EDIT: This file was generated by vugu. Please regenerate instead of editing or add additional code in a separate file.

import "fmt"
import "reflect"
import "github.com/vugu/vjson"
import "github.com/vugu/vugu"
import js "github.com/vugu/vugu/js"

import "github.com/powerman/tr/web/app/internal/ui/comp"
import "log"

var _ = log.Println	// Ensure import is used.
func (c *PageReport) Build(vgin *vugu.BuildIn) (vgout *vugu.BuildOut) {

	vgout = &vugu.BuildOut{}

	var vgiterkey interface{}
	_ = vgiterkey
	var vgn *vugu.VGNode
	vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "div", Attr: []vugu.VGAttribute{vugu.VGAttribute{Namespace: "", Key: "class", Val: "ui-page-report"}}}
	vgout.Out = append(vgout.Out, vgn)	// root for output
	{
		vgparent := vgn
		_ = vgparent
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n\n    "}
		vgparent.AppendChild(vgn)
		{
			vgcompKey := vugu.MakeCompKey(0xF39E4FF9EB9F1BEF^vgin.CurrentPositionHash(), vgiterkey)
			// ask BuildEnv for prior instance of this specific component
			vgcomp, _ := vgin.BuildEnv.CachedComponent(vgcompKey).(*comp.RecordForm)
			if vgcomp == nil {
				// create new one if needed
				vgcomp = new(comp.RecordForm)
				vgin.BuildEnv.WireComponent(vgcomp)
			}
			vgin.BuildEnv.UseComponent(vgcompKey, vgcomp)	// ensure we can use this in the cache next time around
			vgout.Components = append(vgout.Components, vgcomp)
			vgn = &vugu.VGNode{Component: vgcomp}
			vgparent.AppendChild(vgn)
		}
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n\n    "}
		vgparent.AppendChild(vgn)
		for key, value := range c.Records {
			var vgiterkey interface{} = key
			_ = vgiterkey
			key := key
			_ = key
			value := value
			_ = value
			{
				vgcompKey := vugu.MakeCompKey(0x1DE24ECFB65DCA56^vgin.CurrentPositionHash(), vgiterkey)
				// ask BuildEnv for prior instance of this specific component
				vgcomp, _ := vgin.BuildEnv.CachedComponent(vgcompKey).(*comp.RecordRow)
				if vgcomp == nil {
					// create new one if needed
					vgcomp = new(comp.RecordRow)
					vgin.BuildEnv.WireComponent(vgcomp)
				}
				vgin.BuildEnv.UseComponent(vgcompKey, vgcomp)	// ensure we can use this in the cache next time around
				vgcomp.Record = *value
				vgout.Components = append(vgout.Components, vgcomp)
				vgn = &vugu.VGNode{Component: vgcomp}
				vgparent.AppendChild(vgn)
			}
		}
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n\n"}
		vgparent.AppendChild(vgn)
	}
	return vgout
}

// 'fix' unused imports
var _ fmt.Stringer
var _ reflect.Type
var _ vjson.RawMessage
var _ js.Value