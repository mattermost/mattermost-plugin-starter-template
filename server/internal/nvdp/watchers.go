package nvdp

import (
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/illbjorn/mm-nvd/nvd"
)

// WatcherSubscribe initiates the process of announcing CVEs to the channel the
// command was invoked in.
func (p *Plugin) WatcherSubscribe(args *model.CommandArgs) {
	// don't register the same channel twice
	if p.cveWG.Worker(args.ChannelId) == nil {
		// register a callback to fire a message to this channel when matching
		// criteria CVEs are identified
		p.cveWG.Register(
			args.ChannelId,
			nvd.NewCVEQuery(),
			func(cve nvd.CVE) {
				// create the post
				post := p.NewCVEPost(cve, args.ChannelId, args.RootId)
				// ship it!
				_, _ = p.API.CreatePost(post)
			})

		// notify of subscription success
		// TODO: this should be an ephemeral message to the user that invoked the
		//       command.
		post := p.NewPost("Subscription successful!", args.ChannelId, args.RootId)
		_, _ = p.API.CreatePost(post)
	}
}

// WatcherUnsubscribe unsubscribes (discontinues CVE announcements) for the
// channel the command was invoked in.
func (p *Plugin) WatcherUnsubscribe(args *model.CommandArgs) {
	// only attempt unsubscribe if we've actually registered a worker for this
	// channel
	if worker := p.cveWG.Worker(args.ChannelId); worker != nil {
		// perform the unregister
		if err := p.cveWG.Unregister(args.ChannelId); err == nil {
			// TODO: really should be doing something with these errors.
			_, _ = p.API.CreatePost(
				p.NewPost(
					"Unsubscribe successful!", args.ChannelId,
					args.RootId))
		} else {
			// TODO: really should be doing something with these errors.
			_, _ = p.API.CreatePost(
				p.NewPost(
					"Unsubscribe failed.",
					args.ChannelId,
					args.RootId))
		}
	}
}

// WatcherSet assigns various query properties associated with the CVE polling.
// Examples of this include nvd.CVETag filter, nvd.CVSSSeverityV2 filtering
// and more.
func (p *Plugin) WatcherSet(
	params map[string]string,
	args *model.CommandArgs,
) {
	// retrieve the worker
	worker := p.cveWG.Worker(args.ChannelId)
	if worker == nil {
		p.API.SendEphemeralPost(
			args.UserId,
			&model.Post{
				Message:   "`/nvdm subscribe` must be called before `/nvdm set`!",
				ChannelId: args.ChannelId,
				RootId:    args.RootId,
			})
		return
	}

	// get the worker's query
	query := worker.CVEQuery()

	// iterate our params map by key-value pair
	for k, v := range params {
		// TODO: there's absolutely a better way to do this. Time constraints!
		switch {
		case strings.EqualFold(k, nvd.QueryKeyCVSSV2Severity):
			// TODO: validate before cast!
			query.CVSSV2Severity(nvd.CVSSSeverityV2(v))
		case strings.EqualFold(k, nvd.QueryKeyKeywordSearch):
			query.KeywordSearch(v)
		case strings.EqualFold(k, nvd.QueryKeyCVETag):
			// TODO: validate before cast!
			query.CVETag(nvd.CVETag(v))
		}
	}
}

// TODO: implement the inverse of WatcherSet.
func (p *Plugin) WatcherUnset(key string, params map[string]string) {
}
