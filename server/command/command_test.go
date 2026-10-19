package command

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest/mock"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mattermost/mattermost-plugin-starter-template/server/channelstatus"
	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore/mocks"
)

func setupTest(t *testing.T) (*plugintest.API, *pluginapi.Client, *mocks.MockKVStore, *channelstatus.Service) {
	api := &plugintest.API{}
	t.Cleanup(func() { api.AssertExpectations(t) })
	client := pluginapi.NewClient(api, &plugintest.Driver{})
	store := mocks.NewMockKVStore(gomock.NewController(t))
	return api, client, store, channelstatus.NewService(api, client, store)
}

func TestNewCommandHandler(t *testing.T) {
	t.Run("returns an error when registration fails", func(t *testing.T) {
		api, client, _, statuses := setupTest(t)
		api.On("RegisterCommand", mock.Anything).Return(errors.New("boom"))

		handler, err := NewCommandHandler(client, statuses)

		assert.Nil(t, handler)
		assert.ErrorContains(t, err, "failed to register /channelstatus command")
	})
}

func TestStatusCommand(t *testing.T) {
	args := func(command string) *model.CommandArgs {
		return &model.CommandArgs{Command: command, UserId: "user-id", ChannelId: "channel-id"}
	}
	newHandler := func(t *testing.T) (Command, *plugintest.API, *mocks.MockKVStore) {
		api, client, store, statuses := setupTest(t)
		api.On("RegisterCommand", mock.MatchedBy(func(cmd *model.Command) bool {
			return cmd.Trigger == statusCommandTrigger
		})).Return(nil)
		handler, err := NewCommandHandler(client, statuses)
		require.NoError(t, err)
		return handler, api, store
	}
	allowEdit := func(api *plugintest.API, allowed bool) {
		api.On("HasPermissionToChannel", "user-id", "channel-id", model.PermissionReadChannel).Return(true)
		api.On("GetChannel", "channel-id").Return(&model.Channel{Id: "channel-id", Type: model.ChannelTypeOpen}, nil)
		api.On("HasPermissionToChannel", "user-id", "channel-id", model.PermissionManagePublicChannelProperties).Return(allowed)
		api.On("LogAuditRec", mock.Anything)
	}

	t.Run("shows the current status", func(t *testing.T) {
		handler, api, store := newHandler(t)
		api.On("HasPermissionToChannel", "user-id", "channel-id", model.PermissionReadChannel).Return(true)
		store.EXPECT().GetChannelStatus("channel-id").Return(channelstatus.AtRisk, nil)

		response, err := handler.Handle(&plugin.Context{}, args("/channelstatus"))

		require.NoError(t, err)
		assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)
		assert.Equal(t, "This channel is **At risk**.", response.Text)
	})

	t.Run("sets the status", func(t *testing.T) {
		handler, api, store := newHandler(t)
		allowEdit(api, true)
		store.EXPECT().SetChannelStatus("channel-id", channelstatus.OffTrack).Return(nil)
		api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything)

		response, err := handler.Handle(&plugin.Context{}, args("/channelstatus off-track"))

		require.NoError(t, err)
		assert.Equal(t, "Channel status set to **Off track**.", response.Text)
	})

	t.Run("clears the status", func(t *testing.T) {
		handler, api, store := newHandler(t)
		allowEdit(api, true)
		store.EXPECT().DeleteChannelStatus("channel-id").Return(nil)
		api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything)

		response, err := handler.Handle(&plugin.Context{}, args("/channelstatus clear"))

		require.NoError(t, err)
		assert.Equal(t, "Channel status cleared.", response.Text)
	})

	t.Run("explains unknown statuses", func(t *testing.T) {
		handler, api, _ := newHandler(t)
		api.On("LogAuditRec", mock.Anything)

		response, err := handler.Handle(&plugin.Context{}, args("/channelstatus shipped"))

		require.NoError(t, err)
		assert.Equal(t, "Unknown status. Use `on-track`, `at-risk`, `off-track` or `clear`.", response.Text)
	})

	t.Run("tells users without permission", func(t *testing.T) {
		handler, api, _ := newHandler(t)
		allowEdit(api, false)

		response, err := handler.Handle(&plugin.Context{}, args("/channelstatus on-track"))

		require.NoError(t, err)
		assert.Equal(t, "You don't have permission to change this channel's status.", response.Text)
	})
}
