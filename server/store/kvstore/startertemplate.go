package kvstore

import (
	"encoding/json"

	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

// greetingCountKeyPrefix namespaces per-channel records. Keys are capped at 150 bytes,
// so keep prefixes short.
const greetingCountKeyPrefix = "greetings_"

type Client struct {
	client *pluginapi.Client
}

func NewKVStore(client *pluginapi.Client) KVStore {
	return Client{
		client: client,
	}
}

func greetingCountKey(channelID string) string {
	return greetingCountKeyPrefix + channelID
}

func (kv Client) GetGreetingCount(channelID string) (int, error) {
	var count int
	if err := kv.client.KV.Get(greetingCountKey(channelID), &count); err != nil {
		return 0, errors.Wrap(err, "failed to get greeting count")
	}
	return count, nil
}

// IncrementGreetingCount uses compare-and-set so concurrent commands, including ones
// handled by other servers in a cluster, never overwrite each other's increments.
func (kv Client) IncrementGreetingCount(channelID string) (int, error) {
	var count int
	err := kv.client.KV.SetAtomicWithRetries(greetingCountKey(channelID), func(oldValue []byte) (any, error) {
		count = 0
		if oldValue != nil {
			if err := json.Unmarshal(oldValue, &count); err != nil {
				return nil, errors.Wrap(err, "failed to decode greeting count")
			}
		}
		count++
		return count, nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment greeting count")
	}
	return count, nil
}

func (kv Client) ResetGreetingCount(channelID string) error {
	return errors.Wrap(kv.client.KV.Delete(greetingCountKey(channelID)), "failed to reset greeting count")
}
