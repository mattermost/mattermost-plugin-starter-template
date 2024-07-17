package nvdp

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/illbjorn/mm-nvd/nvd"
)

// newCVETemplate defines a boilerplate Mattermost message to be sent when
// new CVEs are identified.
const newCVETemplate = `
:warning: **New CVE Identified!** :warning: 
**CVE ID**: %s
**CVSS Base Score**: %.1f
**DESCRIPTION**: %s
**References**
%s
`

// NewCVEPost creates a model.Post using a CVE announcement template and is
// invoked within the callback passed to nvdm.CVEWatcherGroup.Register().
func (p *Plugin) NewCVEPost(cve nvd.CVE, channelID, rootID string) *model.Post {
	// produce description
	description := cve.Key.Descriptions[0].Value

	// assemble URLs
	urls := make([]string, len(cve.Key.References))
	for i := range cve.Key.References {
		urls[i] = cve.Key.References[i].URL
	}
	urlStr := strings.Join(urls, "\n")

	return &model.Post{
		Message: fmt.Sprintf(
			newCVETemplate,
			cve.Key.ID,
			cve.Key.Metrics.CVSSMetricV2[0].CVSSData.BaseScore,
			description,
			urlStr),
		UserId:    p.botID,
		ChannelId: channelID,
		RootId:    rootID,
	}
}

// NewPost creates a boilerplate model.Post for use with p.API.CreatePost, etc.
func (p *Plugin) NewPost(
	message string,
	channelID,
	rootID string,
) *model.Post {
	return &model.Post{
		Message:   message,
		UserId:    p.botID,
		ChannelId: channelID,
		RootId:    rootID,
	}
}
