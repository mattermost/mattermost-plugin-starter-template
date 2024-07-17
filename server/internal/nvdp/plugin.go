package nvdp

import (
	"context"
	"sync"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"

	"github.com/illbjorn/mm-nvd/server/internal/nvdm"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin
	client            *pluginapi.Client
	configurationLock sync.RWMutex
	configuration     *configuration
	// we'll use this context to handle downstream worker cancellation
	ctx context.Context
	// hang onto the user ID of our registered bot for post creation
	botID string
	// the nvdm.CVEWatcherGroup is the watcher daemon manager that'll handle
	// channels subscribing/unsubscribing and modifying configurations related to
	// the announcement of new NVD vulnerabilities.
	cveWG *nvdm.CVEWatcherGroup
}

// OnActivate initializes the CVEWatcherGroup, registers the bot account and
// registers the slash commands.
func (p *Plugin) OnActivate() error {
	// init the plugin api client if necessary
	if p.client == nil {
		// set a top-level context.Context
		p.ctx = context.Background()
		// init the pluginapi client
		p.client = pluginapi.NewClient(p.API, p.Driver)
	}

	// register the bot if necessary
	if p.botID == "" {
		var err error
		p.botID, err = p.client.Bot.EnsureBot(&nvdBot)
		if err != nil {
			return errors.Wrap(err, "failed to ensure bot")
		}
	}

	// init the CVEWatcherGroup if necessary
	if p.cveWG == nil {
		p.cveWG = nvdm.NewCVEWatcherGroup(p.ctx)
	}

	// register slash commands
	if err := p.API.RegisterCommand(cmdNVDM); err != nil {
		return err
	}

	return nil
}
