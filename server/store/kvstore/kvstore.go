package kvstore

// KVStore wraps the plugin's KV records so the rest of the plugin never builds keys or
// decodes values itself. Add a method here for each kind of record the plugin stores.
type KVStore interface {
	// GetGreetingCount returns how many times /hello has been used in a channel.
	GetGreetingCount(channelID string) (int, error)

	// IncrementGreetingCount atomically adds one to a channel's greeting count and returns the new value.
	IncrementGreetingCount(channelID string) (int, error)

	// ResetGreetingCount deletes a channel's greeting count.
	ResetGreetingCount(channelID string) error
}
