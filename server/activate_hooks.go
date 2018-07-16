package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

// OnActivate is invoked when the plugin is activated.
//
// This sample implementation logs a message to the sample channel whenever the plugin is
// activated.
func (p *Plugin) OnActivate() error {
	// It's necessary to do this asynchronously, so as to avoid CreatePost triggering a call
	// to MessageWillBePosted and deadlocking the plugin.
	//
	// See https://mattermost.atlassian.net/browse/MM-11431
	go func() {
		teams, err := p.API.GetTeams()
		if err != nil {
			p.API.LogError(
				"failed to query teams OnActivate",
				"error", err.Error(),
			)
		}

		for _, team := range teams {
			if _, err := p.API.CreatePost(&model.Post{
				UserId:    p.sampleUserId,
				ChannelId: p.sampleChannelIds[team.Id],
				Message: fmt.Sprintf(
					"OnActivate: %s", PluginId,
				),
				Type: "custom_sample_plugin",
				Props: map[string]interface{}{
					"username":     p.Username,
					"channel_name": p.ChannelName,
				},
			}); err != nil {
				p.API.LogError(
					"failed to post OnActivate message",
					"error", err.Error(),
				)
			}

			if err := p.registerCommand(team.Id); err != nil {
				p.API.LogError(
					"failed to register command",
					"error", err.Error(),
				)
			}
		}

	}()

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This sample implementation logs a debug message to the server logs whenever the plugin is
// activated.
func (p *Plugin) OnDeactivate() error {
	// Ideally, we'd post an on deactivate message like in OnActivate, but this is hampered by
	// https://mattermost.atlassian.net/browse/MM-11431?filter=15018
	p.API.LogDebug("OnDeactivate")

	return nil
}
