package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShowHint displays the next progressive hint for the given exercise.
// Hints are stored in exercises/<dir>/hint.md and delimited by "## Hint N" headings.
// Each call reveals one more hint than the previous.
func ShowHint(p *Printer, state *State, manifest *Manifest, exerciseName string, rootDir string) error {
	ex, err := manifest.FindExercise(exerciseName)
	if err != nil {
		return err
	}

	hintPath := filepath.Join(rootDir, "exercises", ex.Dir, "hint.md")
	data, err := os.ReadFile(hintPath)
	if err != nil {
		return fmt.Errorf("no hint.md found for %s at %s", exerciseName, hintPath)
	}

	hints := parseHints(string(data))
	if len(hints) == 0 {
		p.Info("No hints available for " + exerciseName)
		return nil
	}

	used := state.HintsUsed(exerciseName)
	if used >= len(hints) {
		p.Warning(fmt.Sprintf("All %d hints already shown for %s. Try the solution!", len(hints), exerciseName))
		p.Print(hints[len(hints)-1])
		return nil
	}

	state.IncrementHints(exerciseName)
	if err := state.Save(); err != nil {
		return err
	}

	hint := hints[used]
	p.HintHeader(exerciseName, used+1, len(hints))
	p.Print(hint)

	remaining := len(hints) - (used + 1)
	if remaining > 0 {
		p.Info(fmt.Sprintf("%d more hint(s) available. Run 'dragonflylings hint' again if needed.", remaining))
	} else {
		p.Info("That was the last hint. The answer is in your hands.")
	}

	return nil
}

// parseHints splits hint.md content by "## Hint N" headings.
func parseHints(content string) []string {
	var hints []string
	lines := strings.Split(content, "\n")
	var current strings.Builder
	inHint := false

	for _, line := range lines {
		if strings.HasPrefix(line, "## Hint ") {
			if inHint && current.Len() > 0 {
				hints = append(hints, strings.TrimSpace(current.String()))
				current.Reset()
			}
			inHint = true
			current.WriteString(line + "\n")
		} else if inHint {
			current.WriteString(line + "\n")
		}
	}
	if inHint && current.Len() > 0 {
		hints = append(hints, strings.TrimSpace(current.String()))
	}
	return hints
}
