package command

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest/mock"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mattermost/mattermost-plugin-starter-template/server/channelstatus"
	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore/mocks"
)

// The command's behavior is covered against a real server in TestChannelStatusEndpoint.
// Mocks are only for failures a real server can't produce.
func TestNewCommandHandler(t *testing.T) {
	t.Run("returns an error when registration fails", func(t *testing.T) {
		api := &plugintest.API{}
		t.Cleanup(func() { api.AssertExpectations(t) })
		client := pluginapi.NewClient(api, &plugintest.Driver{})
		statuses := channelstatus.NewService(api, client, mocks.NewMockKVStore(gomock.NewController(t)))
		api.On("RegisterCommand", mock.Anything).Return(errors.New("boom"))

		handler, err := NewCommandHandler(client, statuses)

		assert.Nil(t, handler)
		assert.ErrorContains(t, err, "failed to register /channelstatus command")
	})
}
