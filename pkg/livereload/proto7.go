package livereload

const (
	// http://livereload.com/api/protocol/
	proto7    = "http://livereload.com/protocols/official-7"
	cmdHello  = "hello"
	cmdReload = "reload"
)

type (
	msgClientHello struct {
		Command   string   `json:"command"`
		Protocols []string `json:"protocols"`
	}
	msgServerHello struct {
		Command    string   `json:"command"`
		Protocols  []string `json:"protocols"`
		ServerName string   `json:"serverName"`
	}
	msgServerReload struct {
		Command string `json:"command"`
		Path    string `json:"path"`
	}
)

// MsgHello returns LiveReload protocol 7 ServerHello message.
func MsgHello(serverName string) interface{} {
	return msgServerHello{
		Command:    cmdHello,
		Protocols:  []string{proto7},
		ServerName: serverName,
	}
}

// MsgReload returns LiveReload protocol 7 ServerReload message.
func MsgReload(path string) interface{} {
	return msgServerReload{
		Command: cmdReload,
		Path:    path,
	}
}

func validProto(protocols []string) bool {
	for _, proto := range protocols {
		if proto == proto7 {
			return true
		}
	}
	return false
}
