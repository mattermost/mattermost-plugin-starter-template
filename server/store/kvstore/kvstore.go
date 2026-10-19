package kvstore

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
