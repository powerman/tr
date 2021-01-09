package vugudev

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

func StartRebuildWASM(vuguDir string, cmd string, arg ...string) error {
	eventc, err := watch(vuguDir)
	if err != nil {
		return err
	}
	go rebuildWASM(eventc, vuguDir, cmd, arg...)
	return nil
}

func rebuildWASM(eventc <-chan event, vuguDir string, cmd string, arg ...string) {
	runGenerate := startRunner("go", "generate", "./"+vuguDir+"/...")
	runRebuild := startRunner(cmd, arg...)
	for ev := range eventc {
		if ev.Error != nil {
			log.Println("rebuild wasm: event error:", ev.Error)
		} else {
			switch {
			case strings.HasPrefix(filepath.Base(ev.Body.Name), "."):
			case ev.Body.Op == fsnotify.Chmod:
			case strings.HasSuffix(ev.Body.Name, ".vugu"):
				// log.Println("rebuild wasm: event:", ev.Body.Op, ev.Body.Name)
				runGenerate()
			case strings.HasSuffix(ev.Body.Name, ".go"):
				// log.Println("rebuild wasm: event:", ev.Body.Op, ev.Body.Name)
				runRebuild()
			default:
				// log.Println("rebuild wasm: ignore event:", ev.Body.Op, ev.Body.Name)
			}
		}
	}
}

func startRunner(cmd string, arg ...string) func() {
	c := make(chan struct{}, 1)
	go func() {
		const runDelay = time.Millisecond * 200
		t := time.NewTimer(runDelay)
		if !t.Stop() {
			<-t.C
		}
		for {
			select {
			case <-c:
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				t.Reset(runDelay)
			case <-t.C:
				run(cmd, arg...)
			}
		}
	}()
	return func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

func run(cmd string, arg ...string) {
	log.Println("exec:", cmd, strings.Join(arg, " "))
	c := exec.Command(cmd, arg...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	_ = c.Run()
}
