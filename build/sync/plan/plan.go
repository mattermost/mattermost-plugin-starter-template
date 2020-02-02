package plan

import (
	"encoding/json"
	"fmt"
)

// Plan defines the plan for synchronizing a plugin and a template directory.
type Plan struct {
	Checks []Check `json:"checks"`
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
	return nil
}

// Check returns an error if the condition fails.
type Check interface {
	Check(Context) error
}

// Action runs the defined action.
type Action interface {
	Run() error
}

// jsonPlan is used to unmarshal Plan structures.
type jsonPlan struct {
	Checks []struct {
		Type   string          `json:"type"`
		Params json.RawMessage `json:"params"`
	}
}
