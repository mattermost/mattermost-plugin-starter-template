package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const CommandTrigger = "sample_plugin"

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "(true|false)",
		AutoCompleteDesc: "Enables or disables the sample plugin hooks.",
		DisplayName:      "Sample Plugin Command",
		Description:      "A command used to enable or disable the sample plugin hooks.",
	}); err != nil {
		p.API.LogError(
			"failed to register command",
			"error", err.Error(),
		)
	}

	return nil
}

func (p *Plugin) emitStatusChange() {
	p.API.PublishWebSocketEvent("status_change", map[string]interface{}{
		"enabled": !p.disabled,
	}, &model.WebsocketBroadcast{})
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
// API.
//
// This sample implementation responds to a /sample_plugin command, allowing the user to enable
// or disable the sample plugin's hooks functionality (but leave the command and webapp enabled).
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if !strings.HasPrefix(args.Command, "/"+CommandTrigger) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}

	if strings.HasSuffix(args.Command, "true") {
		if !p.disabled {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "The sample plugin hooks are already enabled.",
			}, nil
		}

		p.disabled = false
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Enabled sample plugin hooks.",
		}, nil

	} else if strings.HasSuffix(args.Command, "false") {
		if p.disabled {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "The sample plugin hooks are already disabled.",
			}, nil
		}

		p.disabled = true
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Disabled sample plugin hooks.",
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command action: " + args.Command),
	}, nil
}
