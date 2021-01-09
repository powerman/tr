package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/powerman/must"

	"github.com/powerman/tr/pkg/livereload"
	"github.com/powerman/tr/pkg/vugudev"
)

const (
	lrPath         = "/devserver/livereload"
	staticDir      = "web/static"
	vuguDir        = "web/app"
	rebuildWASMCmd = "scripts/build-web"
)

func main() {
	useLiveReload := flag.Bool("dev", false, "devserver: use livereload")
	addr := flag.String("addr", "127.0.0.1:8844", "listen address")
	flag.Parse()

	mux := http.NewServeMux()
	handleStatic := http.FileServer(http.Dir(staticDir))

	if *useLiveReload {
		lrServer := livereload.NewServer(livereload.ServerConfig{ForceReloadNewClients: true})
		mux.Handle(lrPath, lrServer)

		lrPatch, err := livereload.NewPatch(staticDir, livereload.PatchConfig{Path: lrPath})
		must.NoErr(err)
		handleStatic = lrPatch.AllHTML(handleStatic)

		must.NoErr(vugudev.StartReloadStatic(staticDir, lrServer))
		must.NoErr(vugudev.StartRebuildWASM(vuguDir, rebuildWASMCmd))
	}

	mux.Handle("/", handleStatic)

	log.Printf("Starting HTTP server at %q", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
