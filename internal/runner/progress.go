package runner

import (
	"encoding/json"
	"os"
	"time"
)

const stateFile = ".dragonflylings-state.json"

// ExerciseStatus tracks the state of a single exercise.
type ExerciseStatus struct {
	Status      string    `json:"status"`       // "pending" | "done"
	CompletedAt time.Time `json:"completed_at,omitempty"`
	HintsUsed   int       `json:"hints_used"`
}

// Stats holds aggregate progress statistics.
type Stats struct {
	Total     int       `json:"total"`
	Completed int       `json:"completed"`
	HintsUsed int       `json:"hints_used"`
	StartedAt time.Time `json:"started_at"`
}

// State is the full progress state persisted to .dragonflylings-state.json.
type State struct {
	Version   int                        `json:"version"`
	Exercises map[string]*ExerciseStatus `json:"exercises"`
	Current   string                     `json:"current"`
	Stats     Stats                      `json:"stats"`
}

// LoadState reads the state file from disk, or returns a fresh state if it doesn't exist.
func LoadState(manifest *Manifest) (*State, error) {
	data, err := os.ReadFile(stateFile)
	if os.IsNotExist(err) {
		return initState(manifest), nil
	}
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return initState(manifest), nil
	}
	// Ensure any new exercises from manifest are present in state.
	for _, ex := range manifest.Exercises {
		if _, ok := s.Exercises[ex.Name]; !ok {
			s.Exercises[ex.Name] = &ExerciseStatus{Status: "pending"}
		}
	}
	return &s, nil
}

func initState(manifest *Manifest) *State {
	s := &State{
		Version:   1,
		Exercises: make(map[string]*ExerciseStatus),
		Stats: Stats{
			Total:     len(manifest.Exercises),
			StartedAt: time.Now(),
		},
	}
	for _, ex := range manifest.Exercises {
		s.Exercises[ex.Name] = &ExerciseStatus{Status: "pending"}
	}
	if len(manifest.Exercises) > 0 {
		s.Current = manifest.Exercises[0].Name
	}
	return s
}

// Save writes the state to disk.
func (s *State) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, data, 0644)
}

// MarkDone marks an exercise as completed.
func (s *State) MarkDone(name string) {
	es := s.Exercises[name]
	if es == nil {
		es = &ExerciseStatus{}
		s.Exercises[name] = es
	}
	es.Status = "done"
	es.CompletedAt = time.Now()
	s.recalcStats()
}

// MarkPending resets an exercise to pending.
func (s *State) MarkPending(name string) {
	if es, ok := s.Exercises[name]; ok {
		es.Status = "pending"
		es.CompletedAt = time.Time{}
		es.HintsUsed = 0
	}
	s.recalcStats()
}

// IncrementHints records that a hint was used for the named exercise.
func (s *State) IncrementHints(name string) {
	if es, ok := s.Exercises[name]; ok {
		es.HintsUsed++
	}
	s.recalcStats()
}

// IsDone returns true if the exercise is marked done.
func (s *State) IsDone(name string) bool {
	if es, ok := s.Exercises[name]; ok {
		return es.Status == "done"
	}
	return false
}

// HintsUsed returns how many hints have been used for an exercise.
func (s *State) HintsUsed(name string) int {
	if es, ok := s.Exercises[name]; ok {
		return es.HintsUsed
	}
	return 0
}

// NextPending returns the first exercise not yet done, or empty string if all done.
func (s *State) NextPending(manifest *Manifest) string {
	for _, ex := range manifest.Exercises {
		if !s.IsDone(ex.Name) {
			return ex.Name
		}
	}
	return ""
}

func (s *State) recalcStats() {
	completed := 0
	hints := 0
	for _, es := range s.Exercises {
		if es.Status == "done" {
			completed++
		}
		hints += es.HintsUsed
	}
	s.Stats.Completed = completed
	s.Stats.HintsUsed = hints
}
