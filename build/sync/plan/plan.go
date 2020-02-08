// The plan package handles the synchronization plan.
//
// Each synchronization plan is a set of checks and actions to perform on specified paths
// that will result in the "plugin" repository being updated.
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
		c, err := parseCheck(check.Type, check.Params)
		if err != nil {
			return fmt.Errorf("failed to parse check %q: %w", check.Type, err)
		}
		p.Checks[i] = c
	}

	if len(t.Paths) > 0 {
		p.Paths = make(map[string][]Action)
	}
	for _, path := range t.Paths {
		var err error
		pathActions := make([]Action, len(path.Actions))
		for i, action := range path.Actions {
			var actionConditions []Check
			if len(action.Conditions) > 0 {
				actionConditions = make([]Check, len(action.Conditions))
			}
			for j, check := range action.Conditions {
				actionConditions[j], err = parseCheck(check.Type, check.Params)
				if err != nil {
					return err
				}
			}
			pathActions[i], err = parseAction(action.Type, action.Params, actionConditions)
			if err != nil {
				return err
			}
		}
		p.Paths[path.Path] = pathActions
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

func parseCheck(checkType string, rawParams json.RawMessage) (Check, error) {
	var c Check

	var params interface{}

	switch checkType {
	case "repo_is_clean":
		tc := RepoIsCleanChecker{}
		params = &tc.Params
		c = &tc
	case "exists":
		tc := PathExistsChecker{}
		params = &tc.Params
		c = &tc
	case "file_unaltered":
		tc := FileUnalteredChecker{}
		params = &tc.Params
		c = &tc
	default:
		return nil, fmt.Errorf("unknown checker type %q", checkType)
	}

	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal params for %s: %w", checkType, err)
	}
	return c, nil
}

func parseAction(actionType string, rawParams json.RawMessage, checks []Check) (Action, error) {
	var a Action

	var params interface{}

	switch actionType {
	case "overwrite_file":
		ta := OverwriteFileAction{}
		ta.Conditions = checks
		params = &ta.Params
		a = &ta
	case "overwrite_directory":
		ta := OverwriteDirectoryAction{}
		ta.Conditions = checks
		params = &ta.Params
		a = &ta
	default:
		return nil, fmt.Errorf("unknown action type %q", actionType)
	}

	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal params for %s: %w", actionType, err)
	}
	return a, nil
}
