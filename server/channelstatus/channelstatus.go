// Package channelstatus holds the rules for reading and changing a channel's status. The slash
// command and the HTTP API both go through it, so they enforce the same checks.
package channelstatus

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore"
)

const (
	OnTrack  = "on-track"
	AtRisk   = "at-risk"
	OffTrack = "off-track"
)

// Values lists the statuses in display order.
var Values = []string{OnTrack, AtRisk, OffTrack}

// labels is the fixed set of statuses. Accepting only these keeps users from putting
// arbitrary text in front of everyone in the channel.
var labels = map[string]string{
	OnTrack:  "On track",
	AtRisk:   "At risk",
	OffTrack: "Off track",
}

// Label returns a status's English name, for slash command replies.
func Label(status string) string {
	return labels[status]
}

// WebSocketEvent reaches the webapp as custom_<plugin id>_status_updated.
const WebSocketEvent = "status_updated"

var (
	ErrInvalidStatus = errors.New("unknown status")
	ErrForbidden     = errors.New("not allowed")
)

type Service struct {
	api    plugin.API
	client *pluginapi.Client
	store  kvstore.KVStore
}

func NewService(api plugin.API, client *pluginapi.Client, store kvstore.KVStore) *Service {
	return &Service{api: api, client: client, store: store}
}

// Get returns the channel's status, or "" if none is set.
func (s *Service) Get(userID, channelID string) (string, error) {
	if !s.client.User.HasPermissionToChannel(userID, channelID, model.PermissionReadChannel) {
		return "", ErrForbidden
	}
	return s.store.GetChannelStatus(channelID)
}

// Set changes the channel's status, or clears it when status is "", and tells the channel's
// members. channelID must already be a valid ID. path is the API route or slash command the
// change came from, recorded in the audit log along with denied attempts.
func (s *Service) Set(c *plugin.Context, path, userID, channelID, status string) error {
	if c == nil {
		c = &plugin.Context{}
	}
	auditRec := plugin.MakeAuditRecordWithContext("setChannelStatus", model.AuditStatusFail, c, userID, path)
	defer s.api.LogAuditRec(auditRec)
	model.AddEventParameterToAuditRec(auditRec, "channel_id", channelID)

	if status != "" && labels[status] == "" {
		return ErrInvalidStatus
	}
	model.AddEventParameterToAuditRec(auditRec, "status", status)

	if err := s.checkCanEdit(userID, channelID); err != nil {
		return err
	}

	var err error
	if status == "" {
		err = s.store.DeleteChannelStatus(channelID)
	} else {
		err = s.store.SetChannelStatus(channelID, status)
	}
	if err != nil {
		return err
	}

	// Everyone who receives an event can read its payload, so broadcast to the channel, not to all users.
	s.client.Frontend.PublishWebSocketEvent(WebSocketEvent, map[string]any{
		"channel_id": channelID,
		"status":     status,
	}, &model.WebsocketBroadcast{ChannelId: channelID})

	auditRec.Success()
	return nil
}

// checkCanEdit applies the server's rule for editing a channel header, so admins control
// who can set statuses with the permissions they already manage.
func (s *Service) checkCanEdit(userID, channelID string) error {
	if !s.client.User.HasPermissionToChannel(userID, channelID, model.PermissionReadChannel) {
		return ErrForbidden
	}

	channel, err := s.client.Channel.Get(channelID)
	if err != nil {
		return errors.Wrap(err, "failed to get channel")
	}

	var permission *model.Permission
	switch channel.Type {
	case model.ChannelTypeOpen:
		permission = model.PermissionManagePublicChannelProperties
	case model.ChannelTypePrivate:
		permission = model.PermissionManagePrivateChannelProperties
	default:
		// Any member of a direct or group message can edit it.
		return nil
	}

	if !s.client.User.HasPermissionToChannel(userID, channelID, permission) {
		return ErrForbidden
	}
	return nil
}
