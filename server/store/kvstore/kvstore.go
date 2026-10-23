package kvstore

import (
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

// KVStore wraps the plugin's KV records so the rest of the plugin never builds keys or
// decodes values itself. Add a method here for each kind of record the plugin stores.
type KVStore interface {
	// GetChannelStatus returns a channel's status, or "" if none is set.
	GetChannelStatus(channelID string) (string, error)

	// SetChannelStatus stores a channel's status.
	SetChannelStatus(channelID, status string) error

	// DeleteChannelStatus removes a channel's status.
	DeleteChannelStatus(channelID string) error
}

// channelStatusKeyPrefix namespaces per-channel records. Keys are capped at 150 bytes,
// so keep prefixes short.
const channelStatusKeyPrefix = "status_"

type Client struct {
	client *pluginapi.Client
}

func NewKVStore(client *pluginapi.Client) KVStore {
	return Client{
		client: client,
	}
}

func channelStatusKey(channelID string) string {
	return channelStatusKeyPrefix + channelID
}

func (kv Client) GetChannelStatus(channelID string) (string, error) {
	var status string
	if err := kv.client.KV.Get(channelStatusKey(channelID), &status); err != nil {
		return "", errors.Wrap(err, "failed to get channel status")
	}
	return status, nil
}

func (kv Client) SetChannelStatus(channelID, status string) error {
	_, err := kv.client.KV.Set(channelStatusKey(channelID), status)
	return errors.Wrap(err, "failed to set channel status")
}

func (kv Client) DeleteChannelStatus(channelID string) error {
	return errors.Wrap(kv.client.KV.Delete(channelStatusKey(channelID)), "failed to delete channel status")
}
