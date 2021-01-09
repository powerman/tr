package vugudev

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/purpleidea/mgmt/recwatch"
)

type event = recwatch.Event // Just to avoid dependency on recwatch in the rest of code.

var errNotADir = errors.New("not a directory")

func watch(dir string) (<-chan event, error) {
	if fi, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("%s: %w", dir, err)
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("%s: %w", dir, errNotADir)
	}

	// XXX recwatch detects directory by "/" suffix, so probably it's
	// not compatible with Windows.
	watcher, err := recwatch.NewRecWatcher(dir+"/", true)
	if err != nil {
		return nil, err
	}

	go func() { log.Fatal("failed to watch ", dir, ": ", watcher.Watch()) }()
	return watcher.Events(), nil
}
