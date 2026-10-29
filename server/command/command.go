package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-starter-template/server/channelstatus"
)

type Command interface {
	Handle(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, error)
}

type Handler struct {
	statuses *channelstatus.Service
}

const (
	statusCommandTrigger = "channelstatus"
	clearArgument        = "clear"
)

// NewCommandHandler registers the plugin's slash commands and returns the handler that executes them.
func NewCommandHandler(client *pluginapi.Client, statuses *channelstatus.Service) (Command, error) {
	items := []model.AutocompleteListItem{}
	for _, status := range channelstatus.Values {
		items = append(items, model.AutocompleteListItem{Item: status, HelpText: channelstatus.Label(status)})
	}
	items = append(items, model.AutocompleteListItem{Item: clearArgument, HelpText: "Remove the status"})

	autocomplete := model.NewAutocompleteData(statusCommandTrigger, "[status]", "Show or set this channel's status")
	autocomplete.AddStaticListArgument("Status", false, items)

	err := client.SlashCommand.Register(&model.Command{
		Trigger:          statusCommandTrigger,
		AutoComplete:     true,
		AutoCompleteDesc: "Show or set this channel's status",
		AutoCompleteHint: "[status]",
		AutocompleteData: autocomplete,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to register /%s command", statusCommandTrigger)
	}

	return &Handler{statuses: statuses}, nil
}

// Handle is called from the ExecuteCommand hook for every command registered in NewCommandHandler.
func (h *Handler) Handle(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, error) {
	fields := strings.Fields(args.Command)
	if len(fields) == 0 {
		return ephemeral("Empty command"), nil
	}

	switch trigger := strings.TrimPrefix(fields[0], "/"); trigger {
	case statusCommandTrigger:
		return h.executeStatusCommand(c, args, fields[1:])
	default:
		return ephemeral(fmt.Sprintf("Unknown command: /%s", trigger)), nil
	}
}

func (h *Handler) executeStatusCommand(c *plugin.Context, args *model.CommandArgs, params []string) (*model.CommandResponse, error) {
	if len(params) == 0 {
		status, err := h.statuses.Get(args.UserId, args.ChannelId)
		if err != nil {
			return nil, err
		}
		if status == "" {
			return ephemeral("This channel has no status. Set one with `/channelstatus on-track`, `at-risk` or `off-track`."), nil
		}
		return ephemeral(fmt.Sprintf("This channel is **%s**.", channelstatus.Label(status))), nil
	}

	status := params[0]
	if status == clearArgument {
		status = ""
	}

	err := h.statuses.Set(c, "/"+statusCommandTrigger, args.UserId, args.ChannelId, status)
	switch {
	case errors.Is(err, channelstatus.ErrInvalidStatus):
		return ephemeral("Unknown status. Use `on-track`, `at-risk`, `off-track` or `clear`."), nil
	case errors.Is(err, channelstatus.ErrForbidden):
		return ephemeral("You don't have permission to change this channel's status."), nil
	case err != nil:
		return nil, err
	case status == "":
		return ephemeral("Channel status cleared."), nil
	default:
		return ephemeral(fmt.Sprintf("Channel status set to **%s**.", channelstatus.Label(status))), nil
	}
}

func ephemeral(text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text,
	}
}
