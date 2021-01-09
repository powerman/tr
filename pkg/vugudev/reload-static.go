// Package vugudev provide helpers for running devserver without vgrun.
package vugudev

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

type Reloader interface {
	Reload(path string)
}

func StartReloadStatic(staticDir string, reloader Reloader) error {
	eventc, err := watch(staticDir)
	if err != nil {
		return err
	}
	go reloadStatic(eventc, reloader)
	return nil
}

func reloadStatic(eventc <-chan event, reloader Reloader) {
	reload := startReloader(reloader)
	for ev := range eventc {
		if ev.Error != nil {
			log.Println("reload static: event error:", ev.Error)
		} else {
			switch {
			case strings.HasSuffix(ev.Body.Name, "~"):
			case strings.HasSuffix(ev.Body.Name, ".swp"):
			case strings.HasSuffix(ev.Body.Name, ".swx"):
			case strings.HasSuffix(ev.Body.Name, ".bak"):
			case strings.HasSuffix(ev.Body.Name, "-go-tmp-umask"):
			case strings.HasPrefix(filepath.Base(ev.Body.Name), "."):
			case ev.Body.Op == fsnotify.Chmod:
			default:
				// log.Println("reload static: event:", ev.Body.Op, ev.Body.Name)
				reload <- ev.Body.Name
			}
		}
	}
}

func startReloader(reloader Reloader) chan<- string {
	c := make(chan string)
	go func() {
		names := make(map[string]struct{})
		const runDelay = time.Millisecond * 200
		t := time.NewTimer(runDelay)
		if !t.Stop() {
			<-t.C
		}
		for {
			select {
			case name := <-c:
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				t.Reset(runDelay)
				names[name] = struct{}{}
			case <-t.C:
				for name := range names {
					delete(names, name)
					reloader.Reload(name)
					log.Println("reloading", name)
				}
			}
		}
	}()
	return c
}
