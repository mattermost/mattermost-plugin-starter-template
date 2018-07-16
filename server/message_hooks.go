package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// MessageWillBePosted is invoked when a message is posted by a user before it is committed to the
// database. If you also want to act on edited posts, see MessageWillBeUpdated. Return values
// should be the modified post or nil if rejected and an explanation for the user.
//
// If you don't need to modify or reject posts, use MessageHasBeenPosted instead.
//
// Note that this method will be called for posts created by plugins, including the plugin that created the post.
//
// This sample implementation rejects posts in the sample channel, as well as posts that @-mention
// the sample plugin user.
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	if p.disabled {
		return post, ""
	}

	// Always allow posts by the sample plugin user.
	if post.UserId == p.sampleUserId {
		return post, ""
	}

	// Reject posts by other users in the sample channels, effectively making it read-only.
	for _, channelId := range p.sampleChannelIds {
		if channelId == post.ChannelId {
			p.API.SendEphemeralPost(post.UserId, &model.Post{
				UserId:    p.sampleUserId,
				ChannelId: channelId,
				Message:   "Posting is not allowed in this channel.",
			})

			return nil, "disallowing post in sample channel"
		}
	}

	// Reject posts mentioning the sample plugin user.
	if strings.Contains(post.Message, fmt.Sprintf("@%s", p.Username)) {
		p.API.SendEphemeralPost(post.UserId, &model.Post{
			UserId:    p.sampleUserId,
			ChannelId: post.ChannelId,
			Message:   "You must not talk about the sample plugin user.",
		})

		return nil, "disallowing mention of sample plugin user"
	}

	// Otherwise, allow the post through.
	return post, ""
}

// MessageWillBeUpdated is invoked when a message is updated by a user before it is committed to
// the database. If you also want to act on new posts, see MessageWillBePosted. Return values
// should be the modified post or nil if rejected and an explanation for the user. On rejection,
// the post will be kept in its previous state.
//
// If you don't need to modify or rejected updated posts, use MessageHasBeenUpdated instead.
//
// Note that this method will be called for posts updated by plugins, including the plugin that
// updated the post.
//
// This sample implementation rejects posts that @-mention the sample plugin user.
func (p *Plugin) MessageWillBeUpdated(c *plugin.Context, newPost, oldPost *model.Post) (*model.Post, string) {
	if p.disabled {
		return newPost, ""
	}

	// Reject posts mentioning the sample plugin user.
	if strings.Contains(newPost.Message, fmt.Sprintf("@%s", p.Username)) {
		p.API.SendEphemeralPost(newPost.UserId, &model.Post{
			UserId:    p.sampleUserId,
			ChannelId: newPost.ChannelId,
			Message:   "You must not talk about the sample plugin user.",
		})

		return nil, "disallowing mention of sample plugin user"
	}

	// Otherwise, allow the post through.
	return newPost, ""
}

// MessageHasBeenPosted is invoked after the message has been committed to the database. If you
// need to modify or reject the post, see MessageWillBePosted Note that this method will be called
// for posts created by plugins, including the plugin that created the post.
//
// This sample implementation logs a message to the sample channel whenever a message is posted,
// unless by the sample plugin user itself.
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if p.disabled {
		return
	}

	// Ignore posts by the sample plugin user.
	if post.UserId == p.sampleUserId {
		return
	}

	user, err := p.API.GetUser(post.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", post.UserId)
		return
	}

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", post.ChannelId)
		return
	}

	if _, err := p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[channel.TeamId],
		Message: fmt.Sprintf(
			"MessageHasBeenPosted in ~%s by @%s",
			channel.Name,
			user.Username,
		),
	}); err != nil {
		p.API.LogError(
			"failed to post MessageHasBeenPosted message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}

// MessageHasBeenUpdated is invoked after a message is updated and has been updated in the
// database. If you need to modify or reject the post, see MessageWillBeUpdated Note that this
// method will be called for posts created by plugins, including the plugin that created the post.
//
// This sample implementation logs a message to the sample channel whenever a message is updated,
// unless by the sample plugin user itself.
func (p *Plugin) MessageHasBeenUpdated(c *plugin.Context, newPost, oldPost *model.Post) {
	if p.disabled {
		return
	}

	// Ignore updates by the sample plugin user.
	if newPost.UserId == p.sampleUserId {
		return
	}

	user, err := p.API.GetUser(newPost.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", newPost.UserId)
		return
	}

	channel, err := p.API.GetChannel(newPost.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", newPost.ChannelId)
		return
	}

	if _, err := p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[channel.TeamId],
		Message: fmt.Sprintf(
			"MessageHasBeenUpdated in ~%s by @%s",
			channel.Name,
			user.Username,
		),
	}); err != nil {
		p.API.LogError(
			"failed to post MessageHasBeenUpdated message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}
