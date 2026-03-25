package runner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch starts a filesystem watcher on the exercises directory.
// When any .go file changes, it re-runs the current exercise.
func Watch(r *Runner) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot create watcher: %w", err)
	}
	defer watcher.Close()

	exercisesDir := filepath.Join(r.RootDir, "exercises")
	if err := addRecursive(watcher, exercisesDir); err != nil {
		return fmt.Errorf("cannot watch exercises dir: %w", err)
	}

	r.Printer.Info("Watch mode active. Watching exercises/ for changes...")
	r.Printer.Info("Edit a .go file to run the current exercise. Ctrl+C to stop.")
	fmt.Println()

	// Run the current exercise once on start.
	_ = r.RunNext()

	var debounce <-chan time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}
			// Debounce rapid saves.
			debounce = time.After(200 * time.Millisecond)

		case <-debounce:
			debounce = nil
			fmt.Printf("\n[%s] Change detected, re-running...\n", time.Now().Format("15:04:05"))
			name := r.State.Current
			if name == "" {
				name = r.State.NextPending(r.Manifest)
			}
			if name != "" {
				_ = r.RunExercise(name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			r.Printer.Warning("Watcher error: " + err.Error())
		}
	}
}

// addRecursive adds the given directory and all subdirectories to the watcher.
func addRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return w.Add(path)
		}
		return nil
	})
}
