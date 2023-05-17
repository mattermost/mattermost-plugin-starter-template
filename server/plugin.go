package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/server/public/model"
	"github.com/mattermost/mattermost-server/server/public/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) MessageHasBeenDeleted(c *plugin.Context, deletedPost *model.Post) {
	user, err := p.API.GetUser(deletedPost.UserId)
	if err != nil {
		p.API.LogError(
			"Failed to query user",
			"user_id", deletedPost.UserId,
			"error", err.Error(),
		)
		return
	}

	channel, err := p.API.GetChannel(deletedPost.ChannelId)
	if err != nil {
		p.API.LogError(
			"Failed to query channel",
			"channel_id", deletedPost.ChannelId,
			"error", err.Error(),
		)
		return
	}

	msg := fmt.Sprintf("MessageHasBeenDeleted: @%s, ~%s\n\n%s", user.Username, channel.Name, deletedPost.Message)
	if _, err := p.API.CreatePost(&model.Post{
		ChannelId: channel.Id,
		Message:   msg,
	}); err != nil {
		p.API.LogError(
			"Failed to post MessageHasBeenDeleted message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}
