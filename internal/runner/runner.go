package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runner executes exercises by running go test or go build.
type Runner struct {
	RootDir  string
	Manifest *Manifest
	State    *State
	Printer  *Printer
}

// NewRunner creates a Runner. rootDir is the project root (containing exercises.toml).
func NewRunner(rootDir string) (*Runner, error) {
	manifest, err := LoadManifest(rootDir)
	if err != nil {
		return nil, err
	}
	state, err := LoadState(manifest)
	if err != nil {
		return nil, err
	}
	return &Runner{
		RootDir:  rootDir,
		Manifest: manifest,
		State:    state,
		Printer:  NewPrinter(),
	}, nil
}

// RunNext runs the next incomplete exercise.
func (r *Runner) RunNext() error {
	name := r.State.NextPending(r.Manifest)
	if name == "" {
		r.Printer.Success("All exercises complete! 🎉")
		r.Printer.Progress(r.Manifest, r.State)
		return nil
	}
	return r.RunExercise(name)
}

// RunExercise runs the named exercise and updates state.
func (r *Runner) RunExercise(name string) error {
	ex, err := r.Manifest.FindExercise(name)
	if err != nil {
		return err
	}

	r.State.Current = name
	_ = r.State.Save()

	exerciseDir := filepath.Join(r.RootDir, "exercises", ex.Dir)
	if _, err := os.Stat(exerciseDir); os.IsNotExist(err) {
		return fmt.Errorf("exercise directory not found: %s", exerciseDir)
	}

	fmt.Printf("\nRunning %s...\n", name)

	// Step 1: compile check (exercises are package main without main() — use vet)
	buildOut, err := r.runCmd(exerciseDir, "go", "vet", "./...")
	if err != nil {
		r.Printer.CompileError(name, buildOut)
		return nil
	}

	// Step 2: test check if mode=test
	if ex.Mode == "test" {
		testOut, err := r.runCmd(exerciseDir, "go", "test", "./...", "-v", "-count=1", "-timeout", "30s")
		if err != nil {
			r.Printer.Failure(name, testOut)
			return nil
		}
		r.Printer.Success(name)
		r.State.MarkDone(name)
		_ = r.State.Save()
		r.Printer.FeynmanPrompt(name, ex.Feynman)

		// Advance to next exercise
		next := r.State.NextPending(r.Manifest)
		if next != "" {
			r.State.Current = next
			_ = r.State.Save()
			r.Printer.Info(fmt.Sprintf("Next up: %s", next))
		}
		return nil
	}

	// compile-only mode
	r.Printer.Success(name)
	r.State.MarkDone(name)
	_ = r.State.Save()
	r.Printer.FeynmanPrompt(name, ex.Feynman)

	next := r.State.NextPending(r.Manifest)
	if next != "" {
		r.State.Current = next
		_ = r.State.Save()
		r.Printer.Info(fmt.Sprintf("Next up: %s", next))
	}
	return nil
}

// VerifyAll re-runs all completed exercises and reports any regressions.
func (r *Runner) VerifyAll() error {
	r.Printer.Info("Verifying all completed exercises...")
	failures := 0
	for _, ex := range r.Manifest.Exercises {
		if !r.State.IsDone(ex.Name) {
			continue
		}
		exerciseDir := filepath.Join(r.RootDir, "exercises", ex.Dir)
		_, buildErr := r.runCmd(exerciseDir, "go", "vet", "./...")
		if buildErr != nil {
			r.Printer.CompileError(ex.Name, "vet failed during verify")
			failures++
			continue
		}
		if ex.Mode == "test" {
			out, testErr := r.runCmd(exerciseDir, "go", "test", "./...", "-count=1", "-timeout", "30s")
			if testErr != nil {
				r.Printer.Failure(ex.Name, out)
				failures++
			} else {
				colorSuccess.Printf("  ✓ %s\n", ex.Name)
			}
		}
	}
	if failures == 0 {
		r.Printer.Info("All completed exercises still pass.")
	} else {
		return fmt.Errorf("%d exercise(s) regressed", failures)
	}
	return nil
}

// ResetAll resets progress for all exercises.
func (r *Runner) ResetAll() error {
	for _, ex := range r.Manifest.Exercises {
		r.State.MarkPending(ex.Name)
	}
	if len(r.Manifest.Exercises) > 0 {
		r.State.Current = r.Manifest.Exercises[0].Name
	}
	return r.State.Save()
}

// ResetExercise resets progress for a single exercise.
func (r *Runner) ResetExercise(name string) error {
	if _, err := r.Manifest.FindExercise(name); err != nil {
		return err
	}
	r.State.MarkPending(name)
	return r.State.Save()
}

// runCmd runs a command in the given directory and returns combined stdout+stderr output.
func (r *Runner) runCmd(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	return output, err
}
