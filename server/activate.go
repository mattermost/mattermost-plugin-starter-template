package main

import (
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

// OnActivate is executed when the plugin is activated
func (p *Plugin) OnActivate() error {
	p.API.LogDebug("Activating plugin")

	isCompatible, err := p.Helpers.CheckRequiredConfig(manifest.RequiredConfig, p.API.GetConfig())
	if err != nil {
		return err
	}

	if !isCompatible {
		errMsg := fmt.Sprintf("Not activating plugin because it is not compatible with the system. Requirements: %s", model.ConfigToJsonWithoutEmptyFields(manifest.RequiredConfig))
		p.API.LogError(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
