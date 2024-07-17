package nvdp

import (
	"github.com/mattermost/mattermost/server/public/model"
)

const (
	// username for the created bot account
	botUsername = "nvdbot"
	// display name used by the created bot account
	botDisplayName = "NVDBot"
	// description for the bot account we'll register
	botDescription = "Bot account responsible for monitoring and announcing" +
		"assigned Nist Vulnerability Database (NVD) CVE announcements."
)

// nvdBot describes the Mattermost bot account which will respond to slash
// commands and announce new CVEs with matching criteria.
var nvdBot = model.Bot{
	Username:    botUsername,
	DisplayName: botDisplayName,
	Description: botDescription,
}
