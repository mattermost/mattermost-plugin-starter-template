package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore"
)

type Command interface {
	Handle(args *model.CommandArgs) (*model.CommandResponse, error)
}

type Handler struct {
	client  *pluginapi.Client
	kvstore kvstore.KVStore
}

const helloCommandTrigger = "hello"

// NewCommandHandler registers the plugin's slash commands and returns the handler that executes them.
func NewCommandHandler(client *pluginapi.Client, store kvstore.KVStore) (Command, error) {
	err := client.SlashCommand.Register(&model.Command{
		Trigger:          helloCommandTrigger,
		AutoComplete:     true,
		AutoCompleteDesc: "Say hello to someone",
		AutoCompleteHint: "[@username]",
		AutocompleteData: model.NewAutocompleteData(helloCommandTrigger, "[@username]", "Username to say hello to"),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to register /%s command", helloCommandTrigger)
	}

	return &Handler{
		client:  client,
		kvstore: store,
	}, nil
}

// Handle is called from the ExecuteCommand hook for every command registered in NewCommandHandler.
func (c *Handler) Handle(args *model.CommandArgs) (*model.CommandResponse, error) {
	fields := strings.Fields(args.Command)
	if len(fields) == 0 {
		return ephemeral("Empty command"), nil
	}

	switch trigger := strings.TrimPrefix(fields[0], "/"); trigger {
	case helloCommandTrigger:
		return c.executeHelloCommand(args, fields[1:])
	default:
		return ephemeral(fmt.Sprintf("Unknown command: /%s", trigger)), nil
	}
}

func (c *Handler) executeHelloCommand(args *model.CommandArgs, params []string) (*model.CommandResponse, error) {
	if len(params) == 0 {
		return ephemeral("Please specify a username"), nil
	}

	// The server only runs commands in channels the user can post in, so args.ChannelId needs no extra check.
	count, err := c.kvstore.IncrementGreetingCount(args.ChannelId)
	if err != nil {
		return nil, err
	}

	return ephemeral(fmt.Sprintf("Hello, %s! That's %d greeting(s) in this channel.", params[0], count)), nil
}

func ephemeral(text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text,
	}
}
