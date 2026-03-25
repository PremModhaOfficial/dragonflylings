package runner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Exercise represents a single exercise entry from exercises.toml.
type Exercise struct {
	Name        string `toml:"name"`
	Dir         string `toml:"dir"`
	Mode        string `toml:"mode"`   // "compile" or "test"
	Module      string `toml:"module"` // e.g. "00_connect"
	HintCount   int    `toml:"hint_count"`
	Feynman     string `toml:"feynman"` // "explain" | "predict" | "limit-test" | "none"
	Description string `toml:"description"`
}

// Manifest holds the ordered list of exercises parsed from exercises.toml.
type Manifest struct {
	Exercises []Exercise `toml:"exercises"`
}

// LoadManifest reads and parses exercises.toml from the project root.
// rootDir should be the directory containing exercises.toml.
func LoadManifest(rootDir string) (*Manifest, error) {
	path := filepath.Join(rootDir, "exercises", "exercises.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read exercises.toml at %s: %w", path, err)
	}
	var m Manifest
	if _, err := toml.Decode(string(data), &m); err != nil {
		return nil, fmt.Errorf("cannot parse exercises.toml: %w", err)
	}
	return &m, nil
}

// FindExercise returns the exercise with the given name, or an error.
func (m *Manifest) FindExercise(name string) (*Exercise, error) {
	for i := range m.Exercises {
		if m.Exercises[i].Name == name {
			return &m.Exercises[i], nil
		}
	}
	return nil, fmt.Errorf("exercise %q not found in manifest", name)
}

// Modules returns a deduplicated ordered list of module names.
func (m *Manifest) Modules() []string {
	seen := make(map[string]bool)
	var out []string
	for _, ex := range m.Exercises {
		if !seen[ex.Module] {
			seen[ex.Module] = true
			out = append(out, ex.Module)
		}
	}
	return out
}

// ExercisesByModule returns all exercises belonging to the given module.
func (m *Manifest) ExercisesByModule(module string) []Exercise {
	var out []Exercise
	for _, ex := range m.Exercises {
		if ex.Module == module {
			out = append(out, ex)
		}
	}
	return out
}
