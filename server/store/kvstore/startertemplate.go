package kvstore

import (
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

// We expose our calls to the KVStore pluginapi methods through this interface for testability and stability.
// This allows us to better control which values are stored with which keys.

type Client struct {
	client *pluginapi.Client
}

func NewKVStore(client *pluginapi.Client) KVStore {
	return Client{
		client: client,
	}
}

// Sample method to get a key-value pair in the KV store
func (kv Client) GetTemplateData(userID string) (string, error) {
	var templateData string
	err := kv.client.KV.Get("template_key-"+userID, &templateData)
	if err != nil {
		return "", errors.Wrap(err, "failed to get template data")
	}
	return templateData, nil
}
