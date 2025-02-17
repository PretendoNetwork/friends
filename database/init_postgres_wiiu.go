package database

import "github.com/PretendoNetwork/friends/globals"

func initPostgresWiiU() {
	var err error

	_, err = Manager.Exec(`CREATE SCHEMA IF NOT EXISTS wiiu`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[Wii U] Postgres schema created")

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.user_data (
		pid integer PRIMARY KEY,
		show_online boolean DEFAULT true,
		show_current_game boolean DEFAULT true,
		block_friend_requests boolean DEFAULT false,
		comment text DEFAULT '',
		comment_changed bigint DEFAULT 0,
		last_online bigint DEFAULT 0
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.network_account_info (
		pid integer PRIMARY KEY,
		unknown1 integer,
		unknown2 integer,
		birthday bigint
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.principal_basic_info (
		pid integer PRIMARY KEY,
		username text,
		unknown integer
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.mii (
		pid integer PRIMARY KEY,
		name text,
		unknown1 integer,
		unknown2 integer,
		data bytea,
		unknown_datetime bigint
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.friendships (
		id bigserial PRIMARY KEY,
		user1_pid integer,
		user2_pid integer,
		date bigint,
		active boolean,
		UNIQUE (user1_pid, user2_pid)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.blocks (
		id bigserial PRIMARY KEY,
		blocker_pid integer,
		blocked_pid integer,
		title_id bigint,
		title_version integer,
		date bigint,
		UNIQUE (blocker_pid, blocked_pid)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS wiiu.friend_requests (
		id bigserial PRIMARY KEY,
		sender_pid integer,
		recipient_pid integer,
		sent_on bigint,
		expires_on bigint,
		message text,
		received boolean,
		accepted boolean,
		denied boolean
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[Wii U] Postgres tables created")
}
