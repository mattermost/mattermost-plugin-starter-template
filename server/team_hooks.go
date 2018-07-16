package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// UserHasJoinedTeam is invoked after the membership has been committed to the database. If
// actor is not nil, the user was added to the team by the actor.
//
// This sample implementation logs a message to the sample channel in the team whenever a user
// joins the team.
func (p *Plugin) UserHasJoinedTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	if p.disabled {
		return
	}

	user, err := p.API.GetUser(teamMember.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", teamMember.UserId)
		return
	}

	if _, err = p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[teamMember.TeamId],
		Message:   fmt.Sprintf("UserHasJoinedTeam: @%s", user.Username),
	}); err != nil {
		p.API.LogError(
			"failed to post UserHasJoinedTeam message",
			"user_id", teamMember.UserId,
			"error", err.Error(),
		)
	}
}

// UserHasLeftTeam is invoked after the membership has been removed from the database. If actor
// is not nil, the user was removed from the team by the actor.
//
// This sample implementation logs a message to the sample channel in the team whenever a user
// leaves the team.
func (p *Plugin) UserHasLeftTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	if p.disabled {
		return
	}

	user, err := p.API.GetUser(teamMember.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", teamMember.UserId)
		return
	}

	if _, err = p.API.CreatePost(&model.Post{
		UserId:    p.sampleUserId,
		ChannelId: p.sampleChannelIds[teamMember.TeamId],
		Message:   fmt.Sprintf("UserHasLeftTeam: @%s", user.Username),
	}); err != nil {
		p.API.LogError(
			"failed to post UserHasLeftTeam message",
			"user_id", teamMember.UserId,
			"error", err.Error(),
		)
	}
}
