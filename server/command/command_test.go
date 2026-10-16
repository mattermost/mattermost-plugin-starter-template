package command

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest/mock"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore/mocks"
)

func setupTest(t *testing.T) (*plugintest.API, *pluginapi.Client, *mocks.MockKVStore) {
	api := &plugintest.API{}
	t.Cleanup(func() { api.AssertExpectations(t) })
	client := pluginapi.NewClient(api, &plugintest.Driver{})
	store := mocks.NewMockKVStore(gomock.NewController(t))
	return api, client, store
}

func TestNewCommandHandler(t *testing.T) {
	t.Run("returns an error when registration fails", func(t *testing.T) {
		api, client, store := setupTest(t)
		api.On("RegisterCommand", mock.Anything).Return(errors.New("boom"))

		handler, err := NewCommandHandler(client, store)

		assert.Nil(t, handler)
		assert.ErrorContains(t, err, "failed to register /hello command")
	})
}

func TestHelloCommand(t *testing.T) {
	newHandler := func(t *testing.T) (Command, *mocks.MockKVStore) {
		api, client, store := setupTest(t)
		api.On("RegisterCommand", mock.MatchedBy(func(cmd *model.Command) bool {
			return cmd.Trigger == helloCommandTrigger
		})).Return(nil)
		handler, err := NewCommandHandler(client, store)
		require.NoError(t, err)
		return handler, store
	}

	t.Run("greets the user and counts the greeting", func(t *testing.T) {
		handler, store := newHandler(t)
		store.EXPECT().IncrementGreetingCount("channel-id").Return(3, nil)

		response, err := handler.Handle(&model.CommandArgs{Command: "/hello @world", ChannelId: "channel-id"})

		require.NoError(t, err)
		assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)
		assert.Equal(t, "Hello, @world! That's 3 greeting(s) in this channel.", response.Text)
	})

	t.Run("asks for a username", func(t *testing.T) {
		handler, _ := newHandler(t)

		response, err := handler.Handle(&model.CommandArgs{Command: "/hello", ChannelId: "channel-id"})

		require.NoError(t, err)
		assert.Equal(t, "Please specify a username", response.Text)
	})

	t.Run("returns store errors", func(t *testing.T) {
		handler, store := newHandler(t)
		store.EXPECT().IncrementGreetingCount("channel-id").Return(0, errors.New("kv unavailable"))

		response, err := handler.Handle(&model.CommandArgs{Command: "/hello @world", ChannelId: "channel-id"})

		assert.Nil(t, response)
		assert.EqualError(t, err, "kv unavailable")
	})
}
