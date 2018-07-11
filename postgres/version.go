package postgres

import "strings"

// Version for a postgres server
type Version string

// IsWalReplayPausedFunctionName returns the name of the function to verify whether the replication
// log is paused according to the postgres version
func (v Version) IsWalReplayPausedFunctionName() string {
	if v.is10() {
		return "pg_is_wal_replay_paused"
	}
	return "pg_is_xlog_replay_paused"
}

// LastWalReceivedLsnFunctionName returns the name of the function that returns the last write-ahead
// log location received and synced to disk by replication according to the postgres version
func (v Version) LastWalReceivedLsnFunctionName() string {
	if v.is10() {
		return "pg_last_wal_receive_lsn"
	}
	return "pg_last_xlog_receive_location"
}

// WalLsnDiffFunctionName returns the name of the function that returns  the difference between two write-ahead
// log locations
func (v Version) WalLsnDiffFunctionName() string {
	if v.is10() {
		return "pg_wal_lsn_diff"
	}
	return "pg_xlog_location_diff"
}

// LastWalReplayedLsnFunctionName returns the name of the function that returns the last write-ahead
// log location replayed during recovery according to the postgres version
func (v Version) LastWalReplayedLsnFunctionName() string {
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
