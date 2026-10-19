package channelstatus

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

	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore/mocks"
)

const (
	userID    = "user-id"
	channelID = "channel-id"
)

func newTestService(t *testing.T) (*Service, *plugintest.API, *mocks.MockKVStore) {
	api := &plugintest.API{}
	t.Cleanup(func() { api.AssertExpectations(t) })
	store := mocks.NewMockKVStore(gomock.NewController(t))
	return NewService(api, pluginapi.NewClient(api, &plugintest.Driver{}), store), api, store
}

func captureAudit(api *plugintest.API) *model.AuditRecord {
	rec := &model.AuditRecord{}
	api.On("LogAuditRec", mock.Anything).Run(func(args mock.Arguments) {
		*rec = *args.Get(0).(*model.AuditRecord)
	}).Once()
	return rec
}

func allowEdit(api *plugintest.API, channelType model.ChannelType, permission *model.Permission) {
	api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(true)
	api.On("GetChannel", channelID).Return(&model.Channel{Id: channelID, Type: channelType}, nil)
	if permission != nil {
		api.On("HasPermissionToChannel", userID, channelID, permission).Return(true)
	}
}

func TestGet(t *testing.T) {
	t.Run("returns the stored status to channel members", func(t *testing.T) {
		service, api, store := newTestService(t)
		api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(true)
		store.EXPECT().GetChannelStatus(channelID).Return(AtRisk, nil)

		status, err := service.Get(userID, channelID)

		require.NoError(t, err)
		assert.Equal(t, AtRisk, status)
	})

	t.Run("rejects users who can't read the channel", func(t *testing.T) {
		service, api, _ := newTestService(t)
		api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(false)

		_, err := service.Get(userID, channelID)

		assert.ErrorIs(t, err, ErrForbidden)
	})
}

func TestSet(t *testing.T) {
	ctx := &plugin.Context{SessionId: "session-id", IPAddress: "10.0.0.1"}

	t.Run("stores the status, notifies the channel and audits the change", func(t *testing.T) {
		service, api, store := newTestService(t)
		allowEdit(api, model.ChannelTypeOpen, model.PermissionManagePublicChannelProperties)
		store.EXPECT().SetChannelStatus(channelID, AtRisk).Return(nil)
		api.On("PublishWebSocketEvent", WebSocketEvent,
			map[string]any{"channel_id": channelID, "status": AtRisk},
			&model.WebsocketBroadcast{ChannelId: channelID},
		).Once()
		audit := captureAudit(api)

		err := service.Set(ctx, "/channelstatus", userID, channelID, AtRisk)

		require.NoError(t, err)
		assert.Equal(t, "setChannelStatus", audit.EventName)
		assert.Equal(t, model.AuditStatusSuccess, audit.Status)
		assert.Equal(t, userID, audit.Actor.UserId)
		assert.Equal(t, "session-id", audit.Actor.SessionId)
		assert.Equal(t, "/channelstatus", audit.Meta[model.AuditKeyAPIPath])
		assert.Equal(t, channelID, audit.EventData.Parameters["channel_id"])
		assert.Equal(t, AtRisk, audit.EventData.Parameters["status"])
	})

	t.Run("clears the status when it's empty", func(t *testing.T) {
		service, api, store := newTestService(t)
		allowEdit(api, model.ChannelTypePrivate, model.PermissionManagePrivateChannelProperties)
		store.EXPECT().DeleteChannelStatus(channelID).Return(nil)
		api.On("PublishWebSocketEvent", WebSocketEvent, map[string]any{"channel_id": channelID, "status": ""}, mock.Anything).Once()
		audit := captureAudit(api)

		require.NoError(t, service.Set(ctx, "/channelstatus", userID, channelID, ""))
		assert.Equal(t, model.AuditStatusSuccess, audit.Status)
	})

	t.Run("lets any member of a direct message set it", func(t *testing.T) {
		service, api, store := newTestService(t)
		allowEdit(api, model.ChannelTypeDirect, nil)
		store.EXPECT().SetChannelStatus(channelID, OnTrack).Return(nil)
		api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Once()
		captureAudit(api)

		assert.NoError(t, service.Set(ctx, "/channelstatus", userID, channelID, OnTrack))
	})

	t.Run("keeps unknown statuses out of the store and the audit record", func(t *testing.T) {
		service, api, _ := newTestService(t)
		audit := captureAudit(api)

		err := service.Set(ctx, "/channelstatus", userID, channelID, "<b>shipped</b>")

		assert.ErrorIs(t, err, ErrInvalidStatus)
		assert.Equal(t, model.AuditStatusFail, audit.Status)
		assert.NotContains(t, audit.EventData.Parameters, "status")
	})

	for name, tc := range map[string]struct {
		channelType model.ChannelType
		permission  *model.Permission
	}{
		"public channel":  {model.ChannelTypeOpen, model.PermissionManagePublicChannelProperties},
		"private channel": {model.ChannelTypePrivate, model.PermissionManagePrivateChannelProperties},
	} {
		t.Run("audits users without permission to edit a "+name, func(t *testing.T) {
			service, api, _ := newTestService(t)
			api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(true)
			api.On("GetChannel", channelID).Return(&model.Channel{Id: channelID, Type: tc.channelType}, nil)
			api.On("HasPermissionToChannel", userID, channelID, tc.permission).Return(false)
			audit := captureAudit(api)

			err := service.Set(ctx, "/channelstatus", userID, channelID, OnTrack)

			assert.ErrorIs(t, err, ErrForbidden)
			assert.Equal(t, model.AuditStatusFail, audit.Status)
		})
	}

	t.Run("rejects non-members before looking up the channel", func(t *testing.T) {
		service, api, _ := newTestService(t)
		api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(false)
		captureAudit(api)

		assert.ErrorIs(t, service.Set(ctx, "/channelstatus", userID, channelID, OnTrack), ErrForbidden)
	})

	t.Run("doesn't notify the channel when the store fails", func(t *testing.T) {
		service, api, store := newTestService(t)
		allowEdit(api, model.ChannelTypeOpen, model.PermissionManagePublicChannelProperties)
		store.EXPECT().SetChannelStatus(channelID, OnTrack).Return(errors.New("kv unavailable"))
		audit := captureAudit(api)

		err := service.Set(ctx, "/channelstatus", userID, channelID, OnTrack)

		assert.EqualError(t, err, "kv unavailable")
		assert.Equal(t, model.AuditStatusFail, audit.Status)
		api.AssertNotCalled(t, "PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything)
	})
}
