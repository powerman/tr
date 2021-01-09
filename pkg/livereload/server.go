// Package livereload implements server for LiveReload protocol 7.
//
// Everything in this package is safe for concurrent use by multiple goroutines.
package livereload

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"

	"github.com/powerman/tr/pkg/broadcast"
)

const defaultServerName = "Go"

// ServerConfig contains Server configuration.
type ServerConfig struct {
	Name                  string // Server name. Default: "Go".
	ForceReloadNewClients bool   // Force reload when client connects for the first time (loose detection using User-Agent).
}

// Server for LiveReload protocol 7.
type Server struct {
	cfg      ServerConfig
	upgrader *websocket.Upgrader
	reload   *broadcast.Topic
	seen     map[string]bool
}

// NewServer creates and returns new Server.
func NewServer(cfg ServerConfig) *Server {
	if cfg.Name == "" {
		cfg.Name = defaultServerName
	}

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := &Server{
		cfg:      cfg,
		upgrader: upgrader,
		reload:   broadcast.NewTopic(),
		seen:     make(map[string]bool),
	}
	return srv
}

// ServeHTTP implements WebSocket LiveReload protocol 7.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := NewConn(ws)

	c.Send() <- MsgHello(srv.cfg.Name)

	if srv.cfg.ForceReloadNewClients {
		if ua := r.Header.Get("User-Agent"); !srv.seen[ua] {
			srv.seen[ua] = true
			c.Send() <- MsgReload("/force-reload.js") // Just a fake name with .js ext.
		}
	}

	srv.reload.Subscribe(c)
	defer srv.reload.Unsubscribe(c)
	c.Wait()
}

// Reload tells all connected LiveReload clients to reload given path.
// CSS and images will be updated without reloading the whole page.
func (srv *Server) Reload(path string) {
	srv.reload.Broadcast(MsgReload(filepath.ToSlash(path)))
}
