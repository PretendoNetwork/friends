package database

import "github.com/PretendoNetwork/friends-secure/globals"

func initPostgresWiiU() {
	var err error

	_, err = Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS wiiu`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[Wii U] Postgres schema created")

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS wiiu.user_data (
		pid integer PRIMARY KEY,
		show_online boolean,
		show_current_game boolean,
		block_friend_requests boolean,
		comment text,
		comment_changed bigint,
		last_online bigint
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS wiiu.friendships (
		id bigserial PRIMARY KEY,
		user1_pid integer,
		user2_pid integer,
		date bigint,
		active boolean
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS wiiu.blocks (
		id bigserial PRIMARY KEY,
		blocker_pid integer,
		blocked_pid integer,
		date bigint,
		active boolean
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS wiiu.friend_requests (
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
