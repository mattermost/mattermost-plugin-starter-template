package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// ChannelHasBeenCreated is invoked after the channel has been committed to the database.
//
// This sample implementation logs a message to the sample channel whenever a channel is created.
func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	if p.disabled {
		return
	}

	if _, err := p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[channel.TeamId],
		Message:   fmt.Sprintf("ChannelHasBeenCreated: ~%s", channel.Name),
	}); err != nil {
		p.API.LogError(
			"failed to post ChannelHasBeenCreated message",
			"channel_id", channel.Id,
			"error", err.Error(),
		)
	}
}

// UserHasJoinedChannel is invoked after the membership has been committed to the database. If
// actor is not nil, the user was invited to the channel by the actor.
//
// This sample implementation logs a message to the sample channel whenever a user joins a channel.
func (p *Plugin) UserHasJoinedChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	if p.disabled {
		return
	}

	user, err := p.API.GetUser(channelMember.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", channelMember.UserId)
		return
	}

	channel, err := p.API.GetChannel(channelMember.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", channelMember.ChannelId)
		return
	}

	if _, err = p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[channel.TeamId],
		Message:   fmt.Sprintf("UserHasJoinedChannel: @%s, ~%s", user.Username, channel.Name),
	}); err != nil {
		p.API.LogError(
			"failed to post UserHasJoinedChannel message",
			"user_id", channelMember.UserId,
			"error", err.Error(),
		)
	}
}

// UserHasLeftChannel is invoked after the membership has been removed from the database. If
// actor is not nil, the user was removed from the channel by the actor.
//
// This sample implementation logs a message to the sample channel whenever a user leaves a
// channel.
func (p *Plugin) UserHasLeftChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	if p.disabled {
		return
	}

	user, err := p.API.GetUser(channelMember.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", channelMember.UserId)
		return
	}

	channel, err := p.API.GetChannel(channelMember.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", channelMember.ChannelId)
		return
	}

	if _, err = p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[channel.TeamId],
		Message:   fmt.Sprintf("UserHasLeftChannel: @%s, ~%s", user.Username, channel.Name),
	}); err != nil {
		p.API.LogError(
			"failed to post UserHasLeftChannel message",
			"user_id", channelMember.UserId,
			"error", err.Error(),
		)
	}
}
