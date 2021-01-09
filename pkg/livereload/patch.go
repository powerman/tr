package livereload

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultSrc         = `//cdn.jsdelivr.net/npm/livereload-js@3.3.1/dist/livereload.js`
	defaultPath        = `/livereload`
	defaultPatchBefore = `</body>`
	patchFmt           = `<script src="%s?%s"></script>%s`
)

var (
	errNotADir     = errors.New("not a directory")
	errPatchBefore = errors.New("HTML does not contains PatchBefore marker")
)

// PatchConfig contains Patch configuration.
type PatchConfig struct {
	Src         string // Default: "//cdn.jsdelivr.net/npm/livereload-js@3.3.1/dist/livereload.js".
	Host        string // Default: autodetect.
	Port        string // Default: autodetect.
	Path        string // Default: "/livereload".
	PatchBefore string // Default: "</body>"
}

// Patch provides HTTP middleware to inject <script src=livereload.js>
// while serving static html files.
type Patch struct {
	dir string
	cfg PatchConfig
}

// NewPatch creates and returns Patch configured to serve static html
// files from staticDir and inject there <script src="{cfg.Src}"> before
// cfg.PatchBefore marker. Injected script will connects to WebSocket
// server at given cfg.Host, cfg.Port and cfg.Path.
func NewPatch(staticDir string, cfg PatchConfig) (*Patch, error) {
	if fi, err := os.Stat(staticDir); err != nil {
		return nil, fmt.Errorf("%s: %w", staticDir, err)
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("%s: %w", staticDir, errNotADir)
	}

	if cfg.Src == "" {
		cfg.Src = defaultSrc
	}
	if cfg.Path == "" {
		cfg.Path = defaultPath
	}
	if cfg.PatchBefore == "" {
		cfg.PatchBefore = defaultPatchBefore
	}

	patch := &Patch{
		dir: staticDir,
		cfg: cfg,
	}
	return patch, nil
}

// AllHTML is an HTTP middleware which will try to handle urls with path
// suffix /, .html and .htm by serving static file with injected
// livereload.js script tag. For suffix / it'll use "index.html".
// If it fails for any reason (e.g. file not exists or it doesn't contains
// cfg.PatchBefore marker) it'll fallback to next handler.
func (p *Patch) AllHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/") {
			path += "index.html"
		}
		if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".htm") {
			var fi os.FileInfo
			var buf []byte
			f, err := os.Open(filepath.Join(p.dir, path))
			if err == nil {
				defer f.Close()
				fi, err = f.Stat()
			}
			if err == nil {
				buf, err = ioutil.ReadAll(f)
			}
			if err == nil {
				host, port := p.addr(r)
				opts := make(url.Values)
				opts.Set("host", host)
				opts.Set("port", port)
				opts.Set("path", p.cfg.Path[1:])
				opts.Set("maxdelay", "1000")

				patch := fmt.Sprintf(patchFmt, p.cfg.Src, opts.Encode(), p.cfg.PatchBefore)
				prevLen := len(buf)
				buf = bytes.Replace(buf, []byte(p.cfg.PatchBefore), []byte(patch), 1)
				if len(buf) == prevLen {
					err = fmt.Errorf("%w %q: %q", errPatchBefore, p.cfg.PatchBefore, f.Name())
				}
			}
			if err == nil {
				http.ServeContent(w, r, filepath.Base(path), fi.ModTime(), bytes.NewReader(buf))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (p *Patch) addr(r *http.Request) (host, port string) {
	host, port, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
		port = "80"
		if r.TLS != nil {
			port = "443"
		}
	}
	if p.cfg.Host != "" {
		host = p.cfg.Host
	}
	if p.cfg.Port != "" {
		port = p.cfg.Port
	}
	return host, port
}
