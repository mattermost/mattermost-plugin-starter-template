package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// OnActivate is executed when the plugin is activated
func (p *Plugin) OnActivate() error {
	p.API.LogDebug("Activating plugin")

	isCompatible, requirements, err := p.API.CheckRequiredConfig()
	if err != nil {
		return errors.Wrap(err, "Error checking plugin compatibility")
	}

	if !isCompatible {
		errMsg := fmt.Sprintf("Not activating plugin because it is not compatible with the system. Requirements: %s", requirements)
		p.API.LogError(errMsg)
		return errors.New(errMsg)
	}

	return nil
}
