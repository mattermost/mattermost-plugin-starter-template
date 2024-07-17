package main

import (
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/illbjorn/mm-nvd/server/internal/nvdp"
)

func main() {
	plugin.ClientMain(&nvdp.Plugin{})
}
