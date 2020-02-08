package plan_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan"
)

func TestUnmarshalPlan(t *testing.T) {
	assert := assert.New(t)
	rawJson := []byte(`
{
  "checks": [
    {"type": "repo_is_clean", "params": {"repo": "template"}}
  ],
  "paths": [
    {
      "path": "abc",
      "actions": [{
        "type": "overwrite_file",
        "params": {"create": true},
        "conditions": [{
          "type": "exists",
          "params": {"repo": "plugin"}
        }]
      }]
    }
  ]
}`)
	var p plan.Plan
	err := json.Unmarshal(rawJson, &p)
	assert.Nil(err)
	expectedCheck := plan.RepoIsCleanChecker{}
	expectedCheck.Params.Repo = "template"

	expectedAction := plan.OverwriteFileAction{}
	expectedAction.Params.Create = true
	expectedActionCheck := plan.PathExistsChecker{}
	expectedActionCheck.Params.Repo = "plugin"
	expectedAction.Conditions = []plan.Check{&expectedActionCheck}
	expected := plan.Plan{
		Checks: []plan.Check{&expectedCheck},
		Paths: map[string][]plan.Action{
			"abc": []plan.Action{
				&expectedAction,
			},
		},
	}
	assert.Equal(expected, p)
}
