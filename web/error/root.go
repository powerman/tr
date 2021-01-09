package main

import (
	"encoding/base64"

	"github.com/vugu/vugu"
)

var errMsg = "dW5rbm93biBlcnJvcg==" // Should be set by `go build` to actual error in Base64.

type Root struct {
	ErrMsg string
}

func (c *Root) Init(ctx vugu.InitCtx) {
	buf, _ := base64.StdEncoding.DecodeString(errMsg)
	c.ErrMsg = string(buf)
}
