// +build tools

package tools

import (
	_ "github.com/cheekybits/genny"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/vugu/vugu/cmd/vugugen"
	_ "gotest.tools/gotestsum"
)
