package database

import "github.com/PretendoNetwork/friends/globals"

func initPostgres3DS() {
	var err error

	_, err = Manager.Exec(`CREATE SCHEMA IF NOT EXISTS "3ds"`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[3DS] Postgres schema created")

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS "3ds".user_data (
		pid integer PRIMARY KEY,
		show_online boolean DEFAULT true,
		show_current_game boolean DEFAULT true,
		comment text DEFAULT '',
		comment_changed bigint DEFAULT 0,
		last_online bigint DEFAULT 0,
		favorite_title bigint DEFAULT 0,
		favorite_title_version integer DEFAULT 0,
		mii_name text DEFAULT '',
		mii_data bytea DEFAULT '',
		mii_changed bigint DEFAULT 0,
		region integer DEFAULT 0,
		area integer DEFAULT 0,
		language integer DEFAULT 0,
		country integer DEFAULT 0
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Manager.Exec(`CREATE TABLE IF NOT EXISTS "3ds".friendships (
		id bigserial PRIMARY KEY,
		user1_pid integer,
		user2_pid integer,
		type integer,
		date bigint,
		UNIQUE (user1_pid, user2_pid)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[3DS] Postgres tables created")
}
