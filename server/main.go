package main

import (
	"github.com/mattermost/mattermost-server/server/public/plugin"
)

func main() {
	plugin.ClientMain(&Plugin{})
}
