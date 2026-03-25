package runner

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Printer handles all terminal output for the runner.
type Printer struct{}

var (
	colorSuccess = color.New(color.FgGreen, color.Bold)
	colorFailure = color.New(color.FgRed, color.Bold)
	colorInfo    = color.New(color.FgCyan)
	colorWarning = color.New(color.FgYellow)
	colorBold    = color.New(color.Bold)
	colorDim     = color.New(color.FgHiBlack)
	colorHint    = color.New(color.FgMagenta)
)

// NewPrinter creates a Printer.
func NewPrinter() *Printer { return &Printer{} }

// Print outputs plain text.
func (p *Printer) Print(s string) { fmt.Println(s) }

// Info outputs a cyan info line.
func (p *Printer) Info(s string) { colorInfo.Println(s) }

// Warning outputs a yellow warning line.
func (p *Printer) Warning(s string) { colorWarning.Println(s) }

// Success prints a success banner for a completed exercise.
func (p *Printer) Success(name string) {
	fmt.Println()
	colorSuccess.Println("  ✓ " + name + " — PASSED")
	fmt.Println()
}

// Failure prints a failure banner with the test output.
func (p *Printer) Failure(name string, output string) {
	fmt.Println()
	colorFailure.Println("  ✗ " + name + " — FAILED")
	fmt.Println()
	fmt.Println(output)
	colorInfo.Println("  Tip: run 'dragonflylings hint' if you're stuck")
	fmt.Println()
}

// CompileError prints a compile error.
func (p *Printer) CompileError(name string, output string) {
	fmt.Println()
	colorFailure.Println("  ✗ " + name + " — BUILD ERROR")
	fmt.Println()
	fmt.Println(output)
	fmt.Println()
}

// HintHeader prints the hint section header.
func (p *Printer) HintHeader(name string, n, total int) {
	fmt.Println()
	colorHint.Printf("  💡 Hint %d/%d for %s\n", n, total, name)
	colorDim.Println("  " + strings.Repeat("─", 50))
	fmt.Println()
}

// FeynmanPrompt prints the post-exercise Feynman challenge.
func (p *Printer) FeynmanPrompt(name string, feynmanType string) {
	fmt.Println()
	colorBold.Println("  ★ FEYNMAN CHALLENGE")
	colorDim.Println("  " + strings.Repeat("─", 50))
	switch feynmanType {
	case "explain":
		colorInfo.Printf("  Read exercises/%s/explain.md\n", name)
		fmt.Println("  Write your explanation in feynman/explanations/" + name + ".md")
	case "predict":
		colorInfo.Println("  Before moving on: what did you predict vs. what you found?")
		fmt.Println("  Write your reflection in feynman/explanations/" + name + ".md")
	case "limit-test":
		colorInfo.Println("  What's the breaking point of what you just built?")
		fmt.Println("  Document your limit-test findings in feynman/explanations/" + name + ".md")
	}
	fmt.Println("  Add your gaps to feynman/gap_notebook.md")
	fmt.Println()
}

// Progress prints the full progress view.
func (p *Printer) Progress(manifest *Manifest, state *State) {
	fmt.Println()
	colorBold.Println("  DRAGONFLYLINGS PROGRESS")
	colorBold.Println("  =======================")
	fmt.Println()

	for _, mod := range manifest.Modules() {
		exercises := manifest.ExercisesByModule(mod)
		total := len(exercises)
		done := 0
		for _, ex := range exercises {
			if state.IsDone(ex.Name) {
				done++
			}
		}

		bar := progressBar(done, total)
		line := fmt.Sprintf("  %-18s %s %d/%d", mod, bar, done, total)
		if done == total {
			colorSuccess.Println(line)
		} else if done > 0 {
			colorInfo.Println(line)
		} else {
			fmt.Println(line)
		}
	}

	fmt.Println()
	pct := 0.0
	if state.Stats.Total > 0 {
		pct = float64(state.Stats.Completed) / float64(state.Stats.Total) * 100
	}
	colorBold.Printf("  Total: %d/%d (%.1f%%) | Hints used: %d\n",
		state.Stats.Completed, state.Stats.Total, pct, state.Stats.HintsUsed)
	if state.Current != "" {
		colorInfo.Printf("  Current: %s\n", state.Current)
	}
	fmt.Println()
}

// List prints all exercises with their status.
func (p *Printer) List(manifest *Manifest, state *State) {
	fmt.Println()
	colorBold.Println("  EXERCISES")
	colorBold.Println("  =========")
	fmt.Println()

	currentMod := ""
	for _, ex := range manifest.Exercises {
		if ex.Module != currentMod {
			currentMod = ex.Module
			fmt.Println()
			colorBold.Printf("  [ %s ]\n", currentMod)
		}
		status := "○"
		if state.IsDone(ex.Name) {
			colorSuccess.Printf("    ✓ %-20s  %s\n", ex.Name, ex.Description)
		} else if ex.Name == state.Current {
			colorInfo.Printf("  → %-20s  %s\n", ex.Name, ex.Description)
		} else {
			colorDim.Printf("    %s %-20s  %s\n", status, ex.Name, ex.Description)
		}
	}
	fmt.Println()
}

func progressBar(done, total int) string {
	if total == 0 {
		return "[     ]"
	}
	width := 5
	filled := (done * width) / total
	bar := strings.Repeat("#", filled) + strings.Repeat(".", width-filled)
	return "[" + bar + "]"
}
