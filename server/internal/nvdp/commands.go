package nvdp

import (
	"regexp"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

const (
	// defines the base command to interact with the NVD monitor plugin's
	// configuration surrounding the NVD feed monitoring & announcement.
	cmdTrigger        = "nvdm"
	cmdDescription    = "Configure NVD feed monitoring."
	cmdDisplayName    = "NVD Management"
	subCmdSubscribe   = "subscribe"
	subCmdUnsubscribe = "unsubscribe"
	subCmdSet         = "set"
	subCmdUnset       = "unset"
)

// cmdNVDM defines the root command which handles all configuration of monitored
// NVD feed data.
var cmdNVDM = &model.Command{
	Trigger:     cmdTrigger,
	Description: cmdDescription,
	DisplayName: cmdDisplayName,
}

// ExecuteCommand performs NVDM slash command execution.
func (p *Plugin) ExecuteCommand(
	_ *plugin.Context,
	args *model.CommandArgs,
) (*model.CommandResponse, *model.AppError) {
	// parse command
	c, err := parseCommand(args.Command)
	if err != nil {
		return nil, model.NewAppError(
			"",
			"Invalid Subcommand",
			nil,
			err.Error(),
			1)
	}

	// perform the action based on the subcommand called by the user
	switch c.Subcommand {
	case subCmdSubscribe:
		p.WatcherSubscribe(args)
	case subCmdUnsubscribe:
		p.WatcherUnsubscribe(args)
	case subCmdSet:
		p.WatcherSet(c.Parameters, args)
	}

	return &model.CommandResponse{}, &model.AppError{}
}

// command is used to parse user slash command input into a format we can use.
type command struct {
	Subcommand string
	Parameters map[string]string
}

// psSubCmdName describes the pattern string to be used with pSubCmdName below.
var psSubCmdName = `(?i)^/nvdm\s*(subscribe|unsubscribe|set|unset)`

// psSubCmdOptions describes the pattern string to be used with pSubCmdOptions
// below.
var psSubCmdOptions = `(\w+)=(\w+)`

// pSubCmdName matches known subcommands.
var pSubCmdName = regexp.MustCompile(psSubCmdName)

// pSubCmdOptions matches all key-value pairs to be applied as configuration
// options to the query for this channel's CVEWatcher.
var pSubCmdOptions = regexp.MustCompile(psSubCmdOptions)

func parseCommand(line string) (*command, error) {
	// init the command
	c := &command{Parameters: make(map[string]string)}

	// confirm we have a valid subcommand
	if !pSubCmdName.MatchString(line) {
		return c, ErrInvalidSubcommand
	}

	// match the subcommand
	subcommandMatches := pSubCmdName.FindAllStringSubmatch(line, -1)

	// get the subcommand
	// lower-case the subcommand for comparison later
	c.Subcommand = strings.ToLower(subcommandMatches[0][1])

	// if the command given was subscribe or unsubscribe - no need to parse key
	// -value pairs
	if c.Subcommand == subCmdSubscribe || c.Subcommand == subCmdUnsubscribe {
		return c, nil
	}

	// match all key-value pairs
	keyValuePairMatches := pSubCmdOptions.FindAllStringSubmatch(line, -1)

	// load all provided configuration key-value pairs
	for _, m := range keyValuePairMatches {
		// lower-case 'key' for comparison later
		key := strings.ToLower(m[1])
		// the key-value does NOT need to be consistently cased as we'll evaluate
		// these with strings.EqualFold()
		c.Parameters[key] = m[2]
	}

	return c, nil
}
