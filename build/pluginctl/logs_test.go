package main

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckOldest(t *testing.T) {
	for name, tc := range map[string]struct {
		logs           []string
		oldest         string
		expectedLogs   []string
		expecetdOldest string
		expectedAllNew bool
	}{
		"nil logs": {
			logs:           nil,
			oldest:         "oldest",
			expectedLogs:   nil,
			expecetdOldest: "oldest",
			expectedAllNew: false,
		},
		"empty logs": {
			logs:           []string{},
			oldest:         "oldest",
			expectedLogs:   nil,
			expecetdOldest: "oldest",
			expectedAllNew: false,
		},
		"no new entries, one old entry": {
			logs:           []string{"old"},
			oldest:         "old",
			expectedLogs:   []string{},
			expecetdOldest: "old",
			expectedAllNew: false,
		},
		"no new entries, multipile old entries": {
			logs:           []string{"old1", "old2", "old3"},
			oldest:         "old3",
			expectedLogs:   []string{},
			expecetdOldest: "old3",
			expectedAllNew: false,
		},
		"one new entry, no old entry": {
			logs:           []string{"new"},
			oldest:         "old",
			expectedLogs:   []string{"new"},
			expecetdOldest: "new",
			expectedAllNew: true,
		},
		"multipile new entries, no old entry": {
			logs:           []string{"new1", "new2", "new3"},
			oldest:         "old",
			expectedLogs:   []string{"new1", "new2", "new3"},
			expecetdOldest: "new3",
			expectedAllNew: true,
		},
		"one new entry, one old entry": {
			logs:           []string{"old", "new"},
			oldest:         "old",
			expectedLogs:   []string{"new"},
			expecetdOldest: "new",
			expectedAllNew: false,
		},
		"one new entry, multipile old entries": {
			logs:           []string{"old1", "old2", "old3", "new"},
			oldest:         "old3",
			expectedLogs:   []string{"new"},
			expecetdOldest: "new",
			expectedAllNew: false,
		},
		"multipile new entries, ultipile old entries": {
			logs:           []string{"old1", "old2", "old3", "new1", "new2", "new3"},
			oldest:         "old3",
			expectedLogs:   []string{"new1", "new2", "new3"},
			expecetdOldest: "new3",
			expectedAllNew: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			logs, oldest, allNew := checkOldest(tc.logs, tc.oldest)

			if allNew != tc.expectedAllNew {
				t.Logf("expected allNew: %v, got %v", tc.expectedAllNew, allNew)
				t.Fail()
			}
			if oldest != tc.expecetdOldest {
				t.Logf("expected oldest: %v, got %v", tc.expecetdOldest, oldest)
				t.Fail()
			}

			compareSlice(t, tc.expectedLogs, logs)
		})
	}
}

func TestFilterLogEntries(t *testing.T) {
	now := time.Now()

	for name, tc := range map[string]struct {
		logs         []string
		pluginID     string
		since        time.Time
		expectedLogs []string
		expecetdErr  bool
	}{
		"nil slice": {
			logs:         nil,
			expectedLogs: nil,
			expecetdErr:  false,
		},
		"empty slice": {
			logs:         []string{},
			expectedLogs: nil,
			expecetdErr:  false,
		},
		"no JSON": {
			logs: []string{
				`{"foo"`,
			},
			expectedLogs: nil,
			expecetdErr:  true,
		},
		"unknown time format": {
			logs: []string{
				`{"message":"foo", "plugin_id": "some.plugin.id", "timestamp": "2023-12-18 10:58:53"}`,
			},
			pluginID:     "some.plugin.id",
			expectedLogs: nil,
			expecetdErr:  true,
		},
		"one matching entry": {
			logs: []string{
				`{"message":"foo", "plugin_id": "some.plugin.id", "timestamp": "2023-12-18 10:58:53.091 +01:00"}`,
			},
			pluginID: "some.plugin.id",
			expectedLogs: []string{
				`{"message":"foo", "plugin_id": "some.plugin.id", "timestamp": "2023-12-18 10:58:53.091 +01:00"}`,
			},
			expecetdErr: false,
		},
		"filter out non plugin entries": {
			logs: []string{
				`{"message":"bar1", "timestamp": "2023-12-18 10:58:52.091 +01:00"}`,
				`{"message":"foo", "plugin_id": "some.plugin.id", "timestamp": "2023-12-18 10:58:53.091 +01:00"}`,
				`{"message":"bar2", "timestamp": "2023-12-18 10:58:54.091 +01:00"}`,
			},
			pluginID: "some.plugin.id",
			expectedLogs: []string{
				`{"message":"foo", "plugin_id": "some.plugin.id", "timestamp": "2023-12-18 10:58:53.091 +01:00"}`,
			},
			expecetdErr: false,
		},
		"filter out old entries": {
			logs: []string{
				fmt.Sprintf(`{"message":"old2", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Add(-2*time.Second).Format(timeStampFormat)),
				fmt.Sprintf(`{"message":"old1", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Add(-1*time.Second).Format(timeStampFormat)),
				fmt.Sprintf(`{"message":"now", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Format(timeStampFormat)),
				fmt.Sprintf(`{"message":"new1", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Add(1*time.Second).Format(timeStampFormat)),
				fmt.Sprintf(`{"message":"new2", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Add(2*time.Second).Format(timeStampFormat)),
			},
			pluginID: "some.plugin.id",
			since:    now,
			expectedLogs: []string{
				fmt.Sprintf(`{"message":"new1", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Add(1*time.Second).Format(timeStampFormat)),
				fmt.Sprintf(`{"message":"new2", "plugin_id": "some.plugin.id", "timestamp": "%s"}`, time.Now().Add(2*time.Second).Format(timeStampFormat)),
			},
			expecetdErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			logs, err := filterLogEntries(tc.logs, tc.pluginID, tc.since)
			if tc.expecetdErr {
				if err == nil {
					t.Logf("expected error, got nil")
					t.Fail()
				}
			} else {
				if err != nil {
					t.Logf("expected no error, got %v", err)
					t.Fail()
				}
			}
			compareSlice(t, tc.expectedLogs, logs)
		})
	}
}

func compareSlice[S ~[]E, E comparable](t *testing.T, expected, got S) {
	if len(expected) != len(got) {
		t.Logf("expected len: %v, got %v", len(expected), len(got))
		t.FailNow()
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			t.Logf("expected [%d]: %v, got %v", i, expected[i], got[i])
			t.Fail()
		}
	}
}
