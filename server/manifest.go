package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

var manifest = struct {
	ID             string
	Version        string
	RequiredConfig *model.Config
}{
	ID:             "com.mattermost.plugin-starter-template",
	Version:        "0.1.0",
	RequiredConfig: model.ConfigFromJson(strings.NewReader(`{"ServiceSettings":{"EnablePostUsernameOverride":true}}`)),
}
