package postgres

import "strings"

// Version for a postgres server
type Version string

// IsLogReplayPausedFunctionName returns the name of the function to verify wether the replication log is paused
// according to the postgres version
func (v Version) IsLogReplayPausedFunctionName() string {
	if v.is10() {
		return "pg_is_wal_replay_paused"
	}
	return "pg_is_xlog_replay_paused"
}

// LastReceivedLsnFunctionName returns the name the of the function that returns the last received LSN
// according to the postgres version
func (v Version) LastReceivedLsnFunctionName() string {
	if v.is10() {
		return "pg_last_wal_receive_lsn"
	}
	return "pg_last_xlog_receive_location"
}

// LastReplayedLsnFunctionName returns the name the of the function that returns the last replayed LSN
// according to the postgres version
func (v Version) LastReplayedLsnFunctionName() string {
	if v.is10() {
		return "pg_last_wal_replay_lsn"
	}
	return "pg_last_xlog_replay_location"

}

func (v Version) is96() bool {
	return strings.HasPrefix(string(v), "9.6.")
}

func (v Version) is10() bool {
	return strings.HasPrefix(string(v), "10.")
}

// Is96Or10 returns whether this is version 9.6.x, 10.x or not
func (v Version) Is96Or10() bool {
	return v.is96() || v.is10()
}
