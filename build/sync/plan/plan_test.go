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
  ]
}`)
	var p plan.Plan
	err := json.Unmarshal(rawJson, &p)
	assert.Nil(err)
	expectedCheck := plan.RepoIsCleanChecker{}
	expectedCheck.Params.Repo = "template"
	expected := plan.Plan{Checks: []plan.Check{&expectedCheck}}
	assert.Equal(expected, p)
}
