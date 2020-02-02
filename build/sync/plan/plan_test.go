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
    {"type": "nil", "params": {"echo": "yay"}}
  ]
}`)
	var p plan.Plan
	err := json.Unmarshal(rawJson, &p)
	assert.Nil(err)
	expectedCheck := plan.NilCheck{}
	expectedCheck.Params.Echo = "yay"
	expected := plan.Plan{Checks: []plan.Check{&expectedCheck}}
	assert.Equal(expected, p)
}
