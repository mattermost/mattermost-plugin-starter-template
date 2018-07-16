package main

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/plugin"
)

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
//
// The Mattermost-User-Id header will be present if (and only if) the request is by an
// authenticated user.
//
// This sample implementation sends back whether or not the plugin hooks are currently enabled. It
// is used by the web app to recover from a network reconnection and synchronize the state of the
// plugin's hooks.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	var response = struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: !p.disabled,
	}

	responseJSON, _ := json.Marshal(response)

	w.Write(responseJSON)
}
