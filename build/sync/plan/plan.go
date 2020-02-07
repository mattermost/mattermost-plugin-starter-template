package plan

import (
	"encoding/json"
	"fmt"
)

// Plan defines the plan for synchronizing a plugin and a template directory.
type Plan struct {
	Checks []Check `json:"checks"`
	Paths  map[string][]Action
}

// UnmarshalJSON implements the `json.Unmarshaler` interface.
func (p *Plan) UnmarshalJSON(raw []byte) error {
	var t jsonPlan
	if err := json.Unmarshal(raw, &t); err != nil {
		return err
	}
	p.Checks = make([]Check, len(t.Checks))
	for i, check := range t.Checks {
		switch check.Type {
		case "repo_is_clean":
			c := RepoIsCleanChecker{}
			err := json.Unmarshal(check.Params, &c.Params)
			if err != nil {
				return fmt.Errorf("failed to unmarshal params for %s: %w", check.Type, err)
			}
			p.Checks[i] = &c
		}
	}

	if len(t.Paths) > 0 {
		p.Paths = make(map[string][]Action)
	}
	for _, path := range t.Paths {
		pathActions := make([]Action, len(path.Actions))
		for i, action := range path.Actions {
			switch action.Type {
			case "overwrite_directory":
				a := OverwriteDirectoryAction{}
				err := json.Unmarshal(action.Params, &a.Params)
				if err != nil {
					return fmt.Errorf("failed to unmarshal params for %s: %w", action.Type, err)
				}
				pathActions[i] = a
			case "overwrite_file":
				a := OverwriteFileAction{}
				err := json.Unmarshal(action.Params, &a.Params)
				if err != nil {
					return fmt.Errorf("failed to unmarshal params for %s: %w", action.Type, err)
				}
				pathActions[i] = a
			}
		}
	}
	return nil
}

// Execute executes the synchronization plan.
func (p *Plan) Execute(c Setup) error {
	for _, check := range p.Checks {
		err := check.Check("", c) // For pre-sync checks, the path is ignored.
		if err != nil {
			return fmt.Errorf("failed check: %w", err)
		}
	}
	return nil
}

// Check returns an error if the condition fails.
type Check interface {
	Check(string, Setup) error
}

// Action runs the defined action.
type Action interface {
	// Run performs the action on the specified path.
	Run(string, Setup) error
	// Check runs checks associated with the action
	// before running it.
	Check(string, Setup) error
}

// jsonPlan is used to unmarshal Plan structures.
type jsonPlan struct {
	Checks []struct {
		Type   string          `json:"type"`
		Params json.RawMessage `json:"params"`
	}
	Paths []struct {
		Path    string `json:"path"`
		Actions []struct {
			Type       string          `json:"type"`
			Params     json.RawMessage `json:"params"`
			Conditions []struct {
				Type   string          `json:"type"`
				Params json.RawMessage `json:"params"`
			}
		}
	}
}
