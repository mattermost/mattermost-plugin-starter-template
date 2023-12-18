package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

const (
	logsPerPage     = 100 // logsPerPage is the number of log entries to fetch per API call
	timeStampFormat = "2006-01-02 15:04:05.000 Z07:00"
)

// logs fetches the latest 500 log entries from Mattermost,
// and prints only the ones related to the plugin to stdout.
func logs(ctx context.Context, client *model.Client4, pluginID string) error {
	for i := 5 - 1; i >= 0; i-- {
		logs, err := fetchLogs(ctx, client, i, pluginID, time.Unix(0, 0))
		if err != nil {
			return fmt.Errorf("failed to fetch log entries: %w", err)
		}

		err = printLogEntries(logs)
		if err != nil {
			return fmt.Errorf("failed to print logs entries: %w", err)
		}
	}

	return nil
}

// watchLogs fetches log entries from Mattermost and print them to stdout.
// It will return without an error when ctx is canceled.
func watchLogs(ctx context.Context, client *model.Client4, pluginID string) error {
	now := time.Now()
	var oldest string

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for {
				logs, err := fetchLogs(ctx, client, 0, pluginID, now)
				if err != nil {
					return fmt.Errorf("failed to fetch log entries: %w", err)
				}

				var allNew bool
				logs, oldest, allNew = checkOldest(logs, oldest)

				err = printLogEntries(logs)
				if err != nil {
					return fmt.Errorf("failed to print logs entries: %w", err)
				}

				if !allNew {
					// No more logs to fetch
					break
				}
			}
		}
	}
}

// checkOldest check a if logs contains new log entries.
// It returns the filtered slice of log entries, the new oldest entry and whether or not a entries were new.
func checkOldest(logs []string, oldest string) ([]string, string, bool) {
	if len(logs) == 0 {
		return nil, oldest, false
	}

	i := slices.Index(logs, oldest)
	switch i {
	case -1:
		// Every log entry is new
		return logs, logs[(len(logs) - 1)], true
	case len(logs) - 1:
		// No new log entries
		return nil, oldest, false
	default:
		// Filter out oldest log entry
		return logs[i+1:], logs[(len(logs) - 1)], false
	}
}

// fetchLogs fetches log entries from Mattermost
// and filters them based on pluginID and timestamp.
func fetchLogs(ctx context.Context, client *model.Client4, page int, pluginID string, since time.Time) ([]string, error) {
	logs, _, err := client.GetLogs(ctx, page, logsPerPage)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs from Mattermost: %w", err)
	}

	logs, err = filterLogEntries(logs, pluginID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to filter log entries: %w", err)
	}

	return logs, nil
}

// filterLogEntries filters a given slice of log entries by pluginID.
// It also filters out any entries which timestamps are older then since.
func filterLogEntries(logs []string, pluginID string, since time.Time) ([]string, error) {
	type logEntry struct {
		PluginID  string `json:"plugin_id"`
		Timestamp string `json:"timestamp"`
	}

	var ret []string

	for _, e := range logs {
		var le logEntry
		err := json.Unmarshal([]byte(e), &le)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal log entry into JSON: %w", err)
		}
		if le.PluginID != pluginID {
			continue
		}

		let, err := time.Parse(timeStampFormat, le.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("unknown timestamp format: %w", err)
		}
		if since.After(let) {
			continue
		}

		// Log entries returned by the API have a newline a prefix.
		// Remove that to make printing consistent.
		e = strings.TrimPrefix(e, "\n")

		ret = append(ret, e)
	}

	return ret, nil
}

// printLogEntries print a slice of log entries to stdout.
func printLogEntries(entries []string) error {
	for _, e := range entries {
		_, err := io.WriteString(os.Stdout, e+"\n")
		if err != nil {
			return fmt.Errorf("failed to write log entry to stdout: %w", err)
		}
	}

	return nil
}
